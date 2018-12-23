package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"sort"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

var (
	className   string                 = ""
	keyType     string                 = ""
	subKeyType  string                 = ""
	subItemType string                 = ""
	structType  string                 = ""
	fields      map[string]interface{} = nil
	format      string                 = ""
)

func doFile(js *simplejson.Json) error {
	var err error
	className, err = js.Get("name").String()
	if err != nil {
		return err
	}
	keyType, err = js.Get("key").String()
	if err != nil {
		return err
	}
	if keyType == "uint" || keyType == "int" {
		return errors.New("no support type: uint or int. please use uint8 int8 uint16 int16 ... etc")
	}
	structType, err = js.Get("type").String()
	if err != nil {
		return err
	}
	format, err = js.Get("struct_format").String()
	if err != nil || format == "" {
		format = "cstruct-go"
	}
	if (format == "protobuf" || format == "gogo") && structType == "1-n" {
		subItemType, err = js.Get("field").String()
		if err != nil {
			fmt.Println("can't find 'field'")
			return err
		}
	} else {
		fields, err = js.Get("fields").Map()
		if err != nil {
			fmt.Println("can't find 'fields'")
			return err
		}
	}
	if err := checkFieldType(); err != nil {
		return err
	}
	if structType == "1-1" {
		return doType11()
	} else if structType == "1-n" {
		subKeyType, err = js.Get("subkey").String()
		if err != nil {
			return err
		}
		if subKeyType == "uint" || subKeyType == "int" {
			return errors.New("no support type: uint or int. please use uint8 int8 uint16 int16 ... etc")
		}
		return doType1n()
	} else {
		return errors.New("type error. type = " + structType)
	}
}

func doType11() error {
	template := template11
	template = strings.Replace(template, "{{packagename}}", *packageName, -1)
	template = strings.Replace(template, "{{classname}}", className, -1)
	template = strings.Replace(template, "{{key_type}}", keyType, -1)
	template = strings.Replace(template, "{{fields_def}}", getFieldsDef(false), -1)
	template = strings.Replace(template, "{{fields_def_db}}", getFieldsDefDB(), -1)
	template = strings.Replace(template, "{{fields_init}}", getFieldsInit(), -1)
	template = strings.Replace(template, "{{func_get}}", getFuncGet(), -1)
	template = strings.Replace(template, "{{func_set}}", getFuncSet(), -1)
	template = strings.Replace(template, "{{fields_save}}", getFuncSave(getFuncStringSave), -1)
	template = strings.Replace(template, "{{fields_save2}}", getFuncSave(getFuncStringSave2), -1)

	if hasStructField() {
		if format == "cstruct-go" {
			template = strings.Replace(template, "{{import_struct_format}}", "cstruct \"github.com/fananchong/cstruct-go\"", -1)
		} else if format == "json" {
			template = strings.Replace(template, "{{import_struct_format}}", "\"encoding/json\"", -1)
		} else if format == "protobuf" {
			template = strings.Replace(template, "{{import_struct_format}}", "\"github.com/golang/protobuf/proto\"", -1)
		} else if format == "gogo" {
			template = strings.Replace(template, "{{import_struct_format}}", "\"github.com/gogo/protobuf/proto\"", -1)
		} else if format != "" {
			panic("unknow format, format =" + format)
		}
	} else {
		template = strings.Replace(template, "{{import_struct_format}}", "", -1)
	}

	if strings.Contains(keyType, "int") {
		template = strings.Replace(template, "{{func_dbkey}}", getFuncDbKeyInt(), -1)
		template = strings.Replace(template, "{{fmt}}", "\"fmt\"", -1)
	} else if keyType == "string" {
		template = strings.Replace(template, "{{func_dbkey}}", getFuncDbKeyStr(), -1)
		template = strings.Replace(template, "{{fmt}}", "", -1)
	} else {
		return errors.New("key type error. type = " + keyType)
	}

	if format == "cstruct-go" {
		template = strings.Replace(template, "{{struct_format}}", "cstruct", -1)
	} else if format == "json" {
		template = strings.Replace(template, "{{struct_format}}", "json", -1)
	} else if format == "protobuf" || format == "gogo" {
		template = strings.Replace(template, "{{struct_format}}", "proto", -1)
	}

	outpath := *outDir + "/" + strings.ToLower(className) + ".go"
	err := ioutil.WriteFile(outpath, []byte(template), 0666)
	if err != nil {
		return err
	}
	return exec_gofmt(outpath)
}

func exec_gofmt(outpath string) error {
	err := exec.Command("gofmt", "-w", outpath).Run()
	if err != nil {
		err = exec.Command("/usr/local/go/bin/gofmt", "-w", outpath).Run()
	}
	if err != nil {
		err = exec.Command("./goimports", "-w", outpath).Run()
	}
	return err
}

var baseType map[string]int = make(map[string]int)

func isBaseType(k string) bool {
	if len(baseType) == 0 {
		baseType["bool"] = 1
		baseType["int"] = 1
		baseType["int8"] = 1
		baseType["int16"] = 1
		baseType["int32"] = 1
		baseType["int64"] = 1
		baseType["uint"] = 1
		baseType["uint8"] = 1
		baseType["uint16"] = 1
		baseType["uint32"] = 1
		baseType["uint64"] = 1
		baseType["float32"] = 1
		baseType["float64"] = 1
		baseType["string"] = 1
		baseType["[]byte"] = 1
	}
	_, ok := baseType[k]
	return ok
}

func hasStructField() bool {
	has := false
	for _, k := range sortFields() {
		v := fields[k].(string)
		if isBaseType(v) == false {
			has = true
			break
		}
	}
	return has
}

func checkFieldType() error {
	for _, k := range sortFields() {
		v := fields[k].(string)
		if v == "uint" || v == "int" {
			return errors.New("no support type: uint or int. please use uint8 int8 uint16 int16 ... etc")
		}
	}
	return nil
}

func getFieldsDef(up bool) string {
	var ret string = ""
	for _, k := range sortFields() {
		v := fields[k].(string)
		if up == false {
			ret = ret + toLower(k) + " " + v + "\n"
		} else {
			ret = ret + toUpper(k) + " " + v + "\n"
		}
	}
	return ret
}

func getFieldsDefDB() string {
	var ret string = ""
	for _, k := range sortFields() {
		v := fields[k].(string)
		k2 := strings.ToUpper(string(k[0])) + string(k[1:])
		if isBaseType(v) == false {
			ret = ret + k2 + " []byte `redis:\"" + strings.ToLower(k2) + "\"`" + "\n"
		} else {
			ret = ret + k2 + " " + v + " `redis:\"" + strings.ToLower(k2) + "\"`" + "\n"
		}
	}
	return ret
}

func getFieldsInit() string {
	var ret string = ""
	for _, k := range sortFields() {
		if ret != "" {
			ret = ret + "\n"
		}
		v := fields[k].(string)
		if isBaseType(v) == false {
			ret = ret + "if err := {{struct_format}}.Unmarshal(data." + toUpper(k) + ", &" + "this." + toLower(k) + "); err != nil {\n"
			ret = ret + "return err\n"
			ret = ret + "}\n"
		} else {
			ret = ret + "this." + toLower(k) + " = data." + toUpper(k)
		}
	}
	return ret
}

func getFuncGet() string {
	var ret string = ""
	for _, k := range sortFields() {
		v := fields[k].(string)
		template := getFuncString
		if isBaseType(v) == false {
			template = getFuncStringForStructFiled
			template = strings.Replace(template, "{{field_name_lower_all}}", strings.ToLower(k), -1)
		}
		template = strings.Replace(template, "{{classname}}", className, -1)
		template = strings.Replace(template, "{{field_type}}", v, -1)
		template = strings.Replace(template, "{{field_name_upper}}", toUpper(k), -1)
		template = strings.Replace(template, "{{field_name_lower}}", toLower(k), -1)
		ret = ret + template + "\n\n"
	}
	return ret
}

func getFuncSet() string {
	var ret string = ""
	for _, k := range sortFields() {
		v := fields[k].(string)
		if isBaseType(v) == false {
			continue
		}
		template := ""
		if v == "string" {
			template = setFuncString_fieldstring
		} else if v == "[]byte" {
			template = setFuncString_fieldbyte
		} else {
			template = setFuncString
		}
		template = strings.Replace(template, "{{classname}}", className, -1)
		template = strings.Replace(template, "{{field_type}}", v, -1)
		template = strings.Replace(template, "{{field_name_upper}}", toUpper(k), -1)
		template = strings.Replace(template, "{{field_name_lower}}", toLower(k), -1)
		template = strings.Replace(template, "{{field_name_lower_all}}", strings.ToLower(k), -1)
		ret = ret + template + "\n\n"
	}
	return ret
}

func getFuncDbKeyInt() string {
	template := dbkeyFuncString_int
	template = strings.Replace(template, "{{classname}}", className, -1)
	template = strings.Replace(template, "{{key_type}}", keyType, -1)
	return template
}

func getFuncSave(temp string) string {
	var ret string = ""
	for _, k := range sortFields() {
		v := fields[k].(string)
		if isBaseType(v) == false {
			template := temp
			template = strings.Replace(template, "{{field_name_lower}}", toLower(k), -1)
			template = strings.Replace(template, "{{field_name_lower_all}}", strings.ToLower(k), -1)
			if len(ret) != 0 {
				ret = ret + "\n"
			}
			ret = ret + template
		}
	}
	return ret
}

func getFuncDbKeyStr() string {
	template := dbkeyFuncString_str
	template = strings.Replace(template, "{{classname}}", className, -1)
	template = strings.Replace(template, "{{key_type}}", keyType, -1)
	return template
}

func toLower(s string) string {
	return strings.ToLower(string(s[0])) + string(s[1:])
}

func toUpper(s string) string {
	return strings.ToUpper(string(s[0])) + string(s[1:])
}

func sortFields() []string {
	var ret []string
	for k, _ := range fields {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

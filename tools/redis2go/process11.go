package main

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"sort"
	"strings"

	"github.com/bitly/go-simplejson"
)

var (
	className  string                 = ""
	keyType    string                 = ""
	sbuKeyType string                 = ""
	structType string                 = ""
	fields     map[string]interface{} = nil
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
	structType, err = js.Get("type").String()
	if err != nil {
		return err
	}
	fields, err = js.Get("fields").Map()
	if err != nil {
		return err
	}
	if structType == "1-1" {
		return doType11()
	} else if structType == "1-n" {
		sbuKeyType, err = js.Get("subkey").String()
		if err != nil {
			return err
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
	template = strings.Replace(template, "{{fields_def}}", getFieldsDef(), -1)
	template = strings.Replace(template, "{{fields_def_db}}", getFieldsDefDB(), -1)
	template = strings.Replace(template, "{{fields_init}}", getFieldsInit(), -1)
	template = strings.Replace(template, "{{func_get}}", getFuncGet(), -1)
	template = strings.Replace(template, "{{func_set}}", getFuncSet(), -1)

	if strings.Contains(keyType, "int") {
		template = strings.Replace(template, "{{func_dbkey}}", getFuncDbKeyInt(), -1)
		template = strings.Replace(template, "{{fmt}}", "\"fmt\"", -1)
	} else if keyType == "string" {
		template = strings.Replace(template, "{{func_dbkey}}", getFuncDbKeyStr(), -1)
		template = strings.Replace(template, "{{fmt}}", "", -1)
	} else {
		return errors.New("key type error. type = " + keyType)
	}

	outpath := *outDir + "/" + className + ".go"
	err := ioutil.WriteFile(outpath, []byte(template), 0666)
	if err != nil {
		return err
	}
	err = exec.Command("gofmt", "-w", outpath).Run()
	return err
}

func getFieldsDef() string {
	var ret string = ""
	for _, k := range sortFields() {
		v := fields[k].(string)
		ret = ret + toLower(k) + " " + v + "\n"
	}
	return ret
}

func getFieldsDefDB() string {
	var ret string = ""
	for _, k := range sortFields() {
		v := fields[k].(string)
		k2 := strings.ToUpper(string(k[0])) + string(k[1:])
		ret = ret + k2 + " " + v + " `redis:\"" + strings.ToLower(k2) + "\"`" + "\n"
	}
	return ret
}

func getFieldsInit() string {
	var ret string = ""
	for _, k := range sortFields() {
		if ret != "" {
			ret = ret + "\n"
		}
		ret = ret + "this." + toLower(k) + " = data." + toUpper(k)
	}
	return ret
}

func getFuncGet() string {
	var ret string = ""
	for _, k := range sortFields() {
		v := fields[k].(string)
		template := getFuncString
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
		template := setFuncString
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

package main

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"
)

func doType1n() error {

	// template1n
	template := template1n
	template = strings.Replace(template, "{{packagename}}", *packageName, -1)
	template = strings.Replace(template, "{{classname}}", className, -1)
	template = strings.Replace(template, "{{key_type}}", keyType, -1)
	template = strings.Replace(template, "{{sub_key_type}}", sbuKeyType, -1)
	template = strings.Replace(template, "{{fields_def}}", getFieldsDef(false), -1)

	if !strings.Contains(sbuKeyType, "int") && sbuKeyType != "string" {
		return errors.New("subkey type error. type = " + sbuKeyType)
	}
	template = strings.Replace(template, "{{conv_subkey}}", getConvSubKey(), -1)
	if strings.Contains(sbuKeyType, "int") {
		template = strings.Replace(template, "{{strconv}}", "\"strconv\"", -1)
	} else {
		template = strings.Replace(template, "{{strconv}}", "", -1)
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

	outpath := *outDir + "/" + className + ".go"
	err := ioutil.WriteFile(outpath, []byte(template), 0666)
	if err != nil {
		return err
	}
	err = exec.Command("gofmt", "-w", outpath).Run()
	if err != nil {
		return err
	}

	// template1n_subitem
	template = template1n_subitem
	template = strings.Replace(template, "{{packagename}}", *packageName, -1)
	template = strings.Replace(template, "{{classname}}", className, -1)
	template = strings.Replace(template, "{{sub_key_type}}", sbuKeyType, -1)
	template = strings.Replace(template, "{{fields_def}}", getFieldsDef(true), -1)
	template = strings.Replace(template, "{{func_get1n}}", getFuncGet1n(), -1)
	template = strings.Replace(template, "{{func_set1n}}", getFuncSet1n(), -1)
	outpath = *outDir + "/" + className + "Item.go"
	err = ioutil.WriteFile(outpath, []byte(template), 0666)
	if err != nil {
		return err
	}
	err = exec.Command("gofmt", "-w", outpath).Run()
	return err
}

func getConvSubKey() string {
	var template = ""
	if strings.Contains(sbuKeyType, "int") {
		template = convSubKeyFuncString_int
		template = strings.Replace(template, "{{sub_key_type}}", sbuKeyType, -1)
	} else if sbuKeyType == "string" {
		template = convSubKeyFuncString_str
	}
	return template
}

func getFuncGet1n() string {
	var ret string = ""
	for _, k := range sortFields() {
		v := fields[k].(string)
		template := get1nFuncString
		if isBaseType(v) == false {
			template = get1nFuncStringForStructFiled
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

func getFuncSet1n() string {
	var ret string = ""
	for _, k := range sortFields() {
		v := fields[k].(string)
		if isBaseType(v) == false {
			continue
		}
		template := set1nFuncString
		template = strings.Replace(template, "{{classname}}", className, -1)
		template = strings.Replace(template, "{{field_type}}", v, -1)
		template = strings.Replace(template, "{{field_name_upper}}", toUpper(k), -1)
		template = strings.Replace(template, "{{field_name_lower}}", toLower(k), -1)
		ret = ret + template + "\n\n"
	}
	return ret
}

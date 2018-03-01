package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"sort"
	"strings"

	"github.com/bitly/go-simplejson"
)

var (
	className  string                 = ""
	keyType    string                 = ""
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
		return doType1n()
	} else {
		return errors.New("type error. type = " + structType)
	}
}

func doType11() error {
	template := template11
	template = strings.Replace(template, "{{packagename}}", *packageName, 1)
	template = strings.Replace(template, "{{classname}}", className, -1)
	template = strings.Replace(template, "{{key_type}}", keyType, -1)
	template = strings.Replace(template, "{{fields_def}}", getFieldsDef(), 1)
	template = strings.Replace(template, "{{func_get}}", getFuncGet(), 1)
	template = strings.Replace(template, "{{func_set}}", "", 1)
	template = strings.Replace(template, "{{func_dbkey}}", "", 1)

	outpath := *outDir + "/" + className + ".go"
	err := ioutil.WriteFile(outpath, []byte(template), 0666)
	if err != nil {
		return err
	}
	fmt.Println("start go fmt file ...")
	err = exec.Command("gofmt", "-w", outpath).Run()
	return err
}

func getFieldsDef() string {
	var ret string = ""
	for _, k := range sortFields() {
		v := fields[k].(string)
		k2 := strings.ToLower(string(k[0])) + string(k[1:])
		ret = ret + k2 + " " + v + "\n"
	}
	return ret
}

func getFuncGet() string {
	var ret string = ""
	for _, k := range sortFields() {
		v := fields[k].(string)
		template := getFuncString
		template = strings.Replace(template, "{{classname}}", className, 1)
		template = strings.Replace(template, "{{field_type}}", v, 1)
		template = strings.Replace(template, "{{field_name_upper}}", toUpper(k), -1)
		template = strings.Replace(template, "{{field_name_lower}}", toLower(k), 1)
		ret = ret + template + "\n\n"
	}
	return ret
}

func doType1n() error {
	return nil
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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitly/go-simplejson"
)

var (
	inDir       = flag.String("input_dir", "", "Json文件目录路径")
	outDir      = flag.String("output_dir", "", "输出目录路径")
	packageName = flag.String("package", "", "生成Go文件的包名")
)

func main() {
	flag.Parse()

	err := filepath.Walk(*inDir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if strings.HasSuffix(f.Name(), ".json") == false {
			return nil
		}

		fmt.Println("start process file =", f.Name(), "...")

		var jsonStr []byte
		jsonStr, err = ioutil.ReadFile(f.Name())
		if err != nil {
			return err
		}
		var js *simplejson.Json
		js, err = simplejson.NewJson(jsonStr)
		if err != nil {
			return err
		}
		err = doFile(js)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

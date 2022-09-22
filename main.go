package main

import (
	"os"
	"fmt"
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/tr8team/gattai/src/gattai-core"
	"text/template"
)

type GattaiFile struct {
    Version string `yaml:"version"`
    TempFolder string `yaml:"temp_folder"`
	EnforceTargets map[string][]string `yaml:"enforce_targets"`
	Repos map[string]interface{} `yaml:"repos"`
	Targets map[string]interface{} `yaml:"targets"`
}

func tpl_fetch(target interface{}) string {
	switch v := target.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
	return "yes"
}

func main() {
	argsWithoutProg := os.Args[1:]

    fmt.Println(argsWithoutProg)

	print.PrintHello()

	temp, err := template.New("./env.gattai.yaml").Funcs(template.FuncMap{
		"fetch": tpl_fetch,
	  }).ParseFiles("./env.gattai.yaml")
	  if err != nil {
		panic(err)
	  }
    // Capture any error
    if err != nil {
        log.Fatalln(err)
    }

	var gattai_file GattaiFile

	yamlFile, err := ioutil.ReadFile("./env.gattai.yaml")
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
	err = yaml.Unmarshal(yamlFile, &gattai_file)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }
    // Print out the template to std
    temp.Execute(os.Stdout, gattai_file.Targets)
}

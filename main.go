package main

import (
	"os"
	"fmt"
	"log"
	//"net/url"
	"io/ioutil"
	//"encoding/json"
	"gopkg.in/yaml.v2"
	"github.com/tr8team/gattai/src/gattai-core"
	"text/template"
)

type Arguments map[string]interface{}

type Target struct {
	Exec string `yaml:"exec"`
	Args Arguments `yaml:"args"`
}

type GattaiFile struct {
    Version string `yaml:"version"`
    TempFolder string `yaml:"temp_folder"`
	EnforceTargets map[string][]string `yaml:"enforce_targets"`
	Repos map[string]map[string]string `yaml:"repos"`
	Targets map[string]map[string]Target `yaml:"targets"`
}

type CLIFile struct {
	Version string `yaml:"version"`

}

func tpl_fetch(lookup map[string]string) func(target Target) string {
	return func(target Target) string {
		yamlKey, err := yaml.Marshal(target)
		if err != nil {
			panic(err)
		}

		result, ok := lookup[string(yamlKey)]
		if !ok {
			result = target.Exec;
			//switch v := target.Args.(type) {
			//case map[interface {}]interface{}:
			//	if exec, ok := v["exec"]; ok {
			//		result = string(exec)
			//	}
			//default:
			//	err := fmt.Sprintf("fetch do not support type %T!\n", v)
			//	panic(err)
			//}

		}
		return result
	}
}



func main() {
	argsWithoutProg := os.Args[1:]

    fmt.Println(argsWithoutProg)

	print.PrintHello()

	var gattai_file GattaiFile

	yamlFile, err := ioutil.ReadFile("env.gattai.yaml")
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
	err = yaml.Unmarshal(yamlFile, &gattai_file)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

	lookUpMap := make(map[string]string)
	//for key, val := gattai_file.Repos {
	//	switch repo, ok := val["repo"];  ok  {
	//	case url.ParseRequestURI(repo):
	//		u, err :=
	//		if err != nil {
	//		   panic(err)
	//		}
	//	}
	//}

	tmpl, err := template.New("env.gattai.yaml").Funcs(template.FuncMap{
		"fetch": tpl_fetch(lookUpMap),
	}).ParseFiles("env.gattai.yaml")
	if err != nil {
		panic(err)
	}
    // Capture any error
    if err != nil {
        log.Fatalln(err)
    }

    // Print out the template to std
	if err := tmpl.Execute(os.Stdout, gattai_file); err != nil {
		fmt.Println(err)
	}
}

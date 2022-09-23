package main

import (
	"os"
	"fmt"
	"log"
	"bytes"
	"strings"
	//"net/url"
	"io/ioutil"
	//"encoding/json"
	"gopkg.in/yaml.v2"
	"github.com/tr8team/gattai/src/gattai-core"
	"text/template"
)

type Target struct {
	Exec string `yaml:"exec"`
	Args map[string]interface{} `yaml:"args"`
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
	Type string `yaml:"type"`
	Params map[string]interface{} `yaml:"params"`
	Return string `yaml:"return"`
	Specs map[string][]map[string]string `yaml:"specs"`
}

func tpl_fetch(tempFolder string, lookUpRepoPath map[string]string) func(target Target) string {
	lookUpReturn := make(map[string]string)
	return func(target Target) string {
		yamlKey, err := yaml.Marshal(target)
		if err != nil {
			panic(err)
		}

		result, ok := lookUpReturn[string(yamlKey)]
		if !ok {
			//result = target.Exec;
			tokens := strings.Split(target.Exec, "/")
			path, ok := lookUpRepoPath[tokens[0]]
			if !ok {

			}
			tmpl_filename := strings.Join(tokens[1:],"/") + ".yaml"
			tmpl_filepath := path + "/" + tmpl_filename
			tmpl, err := template.New(tmpl_filename).Funcs(template.FuncMap{
				"temp_folder": tpl_temp_folder(tempFolder),
			}).ParseFiles(tmpl_filepath)
			if err != nil {
				panic(err)
			}
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, target); err != nil {
				panic(err)
			}
			var cli_file CLIFile;
			err = yaml.Unmarshal(buf.Bytes(), &cli_file)
			if err != nil {
				log.Fatalf("Unmarshal: %v", err)
			}
			result = cli_file.Type
			//switch v := target.Args.(type) {
			//case map[interface {}]interface{}:
			//	if exec, ok := v["exec"]; ok {
			//		result = string(exec)
			//	}
			//default:
			//	err := fmt.Sprintf("fetch do not support type %T!\n", v)
			//	panic(err)
			//}

			// Print out the template to std
		}
		return result
	}
}

func tpl_temp_folder(tempFolder string) func(filename string) string {
	return func(filename string) string {
		return tempFolder + "/" + filename
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

	// TODO: if url.ParseRequestURI(repo):
	// download repo and return path
	lookUpRepoPath := make(map[string]string)
	for key, val := range gattai_file.Repos {
		if repo, ok := val["repo"];  ok  {
			lookUpRepoPath[key] = repo
		}
	}

	tmpl, err := template.New("env.gattai.yaml").Funcs(template.FuncMap{
		"fetch": tpl_fetch(gattai_file.TempFolder,lookUpRepoPath),
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

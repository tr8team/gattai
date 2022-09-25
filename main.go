package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"path"
	"bytes"
	"strings"
	"context"
	"runtime"
	//"os/exec"
	//"net/url"
	"io/ioutil"
	"text/template"
	"gopkg.in/yaml.v2"
	//"mvdan.cc/sh/v3/shell"
	//"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
	"github.com/tr8team/gattai/src/gattai-core"
)

type Target struct {
	Exec string `yaml:"exec"`
	Args map[string]interface{} `yaml:"args"`
}

type Command struct {
	Cmd string `yaml:"cmd"`
	Args []string `yaml:"args"`
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
	Spec map[string][]Command `yaml:"spec"`
}

func tpl_fetch(gattai_file GattaiFile, lookUpRepoPath map[string]string) func(target Target) string {
	lookUpReturn := make(map[string]string)
	return func(target Target) string {
		// get target generated key
		yamlTarget, err := yaml.Marshal(target)
		if err != nil {
			panic(err)
		}
		var buf bytes.Buffer
		// check if result for target already exist
		result, ok := lookUpReturn[string(yamlTarget)]
		if !ok {
			// if not, parse target to see if target have dependency
			tmpl, err := template.New("").Funcs(template.FuncMap{
				"fetch": tpl_fetch(gattai_file,lookUpRepoPath),
			}).Parse(string(yamlTarget))
			if err != nil {
				panic(err)
			}
			buf.Reset()
			if err := tmpl.Execute(&buf, gattai_file); err != nil {
				panic(err)
			}
			// execute return template which hope is the leaf template
			var updated_target Target
			err = yaml.Unmarshal(buf.Bytes(), &updated_target)
			if err != nil {
				log.Fatalf("Unmarshal2: %v", err)
			}
			// unmarshal the update target to create the execution path
			tokens := strings.Split(updated_target.Exec, "/")
			repo_path, ok := lookUpRepoPath[tokens[0]]
			if !ok {

			}
			tmpl_filename := tokens[len(tokens)-1] + ".yaml"
			tmpl_filepath := path.Join(repo_path,path.Join(tokens[1:]...)) + ".yaml"
			tmpl, err = template.New(tmpl_filename).Funcs(template.FuncMap{
				"temp_folder": tpl_temp_folder(gattai_file.TempFolder),
			}).ParseFiles(tmpl_filepath)
			if err != nil {
				panic(err)
			}
			buf.Reset()
			if err := tmpl.Execute(&buf, updated_target); err != nil {
				panic(err)
			}
			var cli_file CLIFile;
			err = yaml.Unmarshal(buf.Bytes(), &cli_file)
			if err != nil {
				log.Fatalf("Unmarshal3: %v", err)
			}
			src := ""
			for _, blk := range cli_file.Spec["cmds"] {

				src += blk.Cmd + " " + strings.Join(blk.Args," ") + ";"
				//expArgs, err := shell.Expand(strings.Join(blk.Args," "),nil)
				//if err != nil {
				//	log.Fatalf("Expand: %v", err)
				//}
				//fmt.Println(expArgs)
				//out, err := shell.Fields(expArgs, nil)
				//if err != nil {
				//	log.Fatalf("Fields: %v", err)
				//}
				//cmd := exec.Command(blk.Cmd,out...)
				//stdout, err := cmd.Output()
				//if err != nil {
				//	log.Fatalf("Command: %v", err)
				//}
				//result = strings.TrimSpace(string(stdout))
				//lookUpReturn[string(yamlTarget)] = result
			}

			file, _ := syntax.NewParser().Parse(strings.NewReader(src), "")

			open := func(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
				if runtime.GOOS == "windows" && path == "/dev/null" {
					path = "NUL"
				}
				return interp.DefaultOpenHandler()(ctx, path, flag, perm)
			}
			buf.Reset()
			runner, _ := interp.New(
				interp.StdIO(nil, &buf, os.Stdout),
				interp.OpenHandler(open),
			)
			err = runner.Run(context.TODO(), file)
			if err != nil {
				log.Fatalf("Run: %v", err)
			}
			result = strings.TrimSpace(buf.String())
			lookUpReturn[string(yamlTarget)] = result
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
        log.Fatalf("Unmarshal1: %v", err)
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
		"fetch": tpl_fetch(gattai_file,lookUpRepoPath),
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

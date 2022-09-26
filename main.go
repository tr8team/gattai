package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"path"
	"time"
	"flag"
	"bytes"
	"errors"
	"strings"
	"context"
	"runtime"
	//"net/url"
	"io/ioutil"
	"text/template"
	"gopkg.in/yaml.v2"
	//"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
	//"github.com/jessevdk/go-flags"
	"github.com/tr8team/gattai/src/gattai-core"
)

type Target struct {
	Exec string `yaml:"exec"`
	Vars map[string]interface{} `yaml:"vars"`
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
	Spec map[string][](map[string]interface{}) `yaml:"spec"`
}

type RunCommand struct {
	fs *flag.FlagSet

	keeptempfiles bool
	destination string
}

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func NewRunCommand() *RunCommand {
	rc := &RunCommand{
		fs: flag.NewFlagSet("run", flag.ContinueOnError),
	}

	rc.fs.BoolVar(&rc.keeptempfiles, "keeptempfiles", false, "Clean up temporary create files")
	rc.fs.StringVar(&rc.destination, "destination", "", "The path where the output will go to")

	return rc
}

func (rc *RunCommand) Name() string {
	return rc.fs.Name()
}

func (rc *RunCommand) Init(args []string) error {
	return rc.fs.Parse(args)
}

func (rc *RunCommand) Run() error {

	args := rc.fs.Args()

	if len(args) < 3 {
		return errors.New("No <namespace> or <target> or <gattai-file> provided!")
	}

	namespace_id := args[0]
	target_id := args[1]
	gattaifile_path := args[2]

	var gattaiFile GattaiFile

	yamlFile, err := ioutil.ReadFile(gattaifile_path)
    if err != nil {
		return fmt.Errorf("Error reading Gattai File: %v", err)
    }
	err = yaml.Unmarshal(yamlFile, &gattaiFile)
    if err != nil {
		return fmt.Errorf("Error parsing Gattai File: %v", err)
    }

	lookUpReturn := make(map[string]string)
	// TODO: if url.ParseRequestURI(repo):
	// download repo and return path
	lookUpRepoPath := make(map[string]string)
	for key, val := range gattaiFile.Repos {
		if repo, ok := val["repo"];  ok  {
			lookUpRepoPath[key] = repo
		}
	}

	tempDir, err := os.MkdirTemp(gattaiFile.TempFolder, "gattai_tmp")
	if err != nil {
		return fmt.Errorf("Error creating temporary folder: %v", err)
	}
	if rc.keeptempfiles == false {
		fmt.Println("Clean up temp files!")
		defer os.RemoveAll(tempDir) // clean up
	}

	switch namespace_id {
	case "*":
		switch  target_id {
		case "*":
			// all namespaces and all targets
			for _, targets := range gattaiFile.Targets {
				for _, target := range targets {
					result := tpl_fetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn)(target)
					fmt.Println(result)
				}
			}
		default:
			// all namespaces and a single target
			for _, targets := range gattaiFile.Targets {
				if target, ok := targets[target_id]; ok {
					result := tpl_fetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn)(target)
					fmt.Println(result)
				}
			}
		}
	default:
		if targets , ok := gattaiFile.Targets[namespace_id]; ok {
			switch  target_id {
			case "*":
				// a single namespace and all targets
				for _, target := range targets {
					result := tpl_fetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn)(target)
					fmt.Println(result)
				}
			default:
				// a single namespace and a single target
				if target, ok := targets[target_id]; ok {
					result := tpl_fetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn)(target)
					fmt.Println(result)
				}
			}
		}
	}

	return nil
}

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command")
	}

	cmds := []Runner{
		NewRunCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

func tpl_fetch(gattai_file GattaiFile, temp_dir string, lookUpRepoPath map[string]string, lookUpReturn map[string]string) func(target Target) string {
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
				"fetch": tpl_fetch(gattai_file,temp_dir,lookUpRepoPath,lookUpReturn),
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
			tmpl_filepath := path.Join(repo_path,path.Join(tokens[1:]...)) + ".yaml"
			tmpl_filename := path.Base(tmpl_filepath)
			tmpl, err = template.New(tmpl_filename).Funcs(template.FuncMap{
				"temp_folder": tpl_temp_folder(temp_dir),
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
				if command, ok := blk["command"].(string); ok {
					src += command
					switch args := blk["args"].(type) {
					case []interface {}:
						for _, elem := range args {
							src +=  " " + elem.(string)
						}
						src +=  ";"
					default:
						err := fmt.Sprintf("fetch do not support type %T!\n", args)
						panic(err)
					}
				}
				//if include, ok := blk["include"].(string); ok {
				//
				//}
			}

			file, _ := syntax.NewParser().Parse(strings.NewReader(src), "")

			open := func(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
				if runtime.GOOS == "windows" && path == "/dev/null" {
					path = "NUL"
				}
				return interp.DefaultOpenHandler()(ctx, path, flag, perm)
			}
			exec := func(ctx context.Context, args []string) error {
				hc := interp.HandlerCtx(ctx)

				//if args[0] == "join" {
				//	fmt.Fprintln(hc.Stdout, strings.Join(args[2:], args[1]))
				//	return nil
				//}

				if _, err := interp.LookPathDir(hc.Dir, hc.Env, args[0]); err != nil {
					fmt.Printf("%s is not installed\n", args[0])
					return interp.NewExitStatus(1)
				}

				return interp.DefaultExecHandler(2*time.Second)(ctx, args)
			}
			buf.Reset()
			runner, _ := interp.New(
				//interp.Env(expand.ListEnviron("GLOBAL=global_value")),
				interp.StdIO(nil, &buf, os.Stdout),
				interp.OpenHandler(open),
				interp.ExecHandler(exec),
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

func tpl_temp_folder(temp_dir string) func(filename string) string {
	return func(filename string) string {
		return path.Join(temp_dir,filename)
	}
}

func main() {
	print.PrintHello()

	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

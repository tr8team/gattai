package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"path"
	"time"
	"bytes"
	"strings"
	"context"
	"runtime"
	"os/exec"
	"net/url"
	"io/ioutil"
	"text/template"
	"gopkg.in/yaml.v2"
	//"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
	"github.com/spf13/cobra"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tr8team/gattai/src/gattai-core"
)

type Target struct {
	Exec string `yaml:"exec"`
	RunTimeEnv string `yaml:"runtime_env"`
	Vars map[string]interface{} `yaml:"vars"`
}

type GattaiFile struct {
    Version string `yaml:"version"`
    TempFolder string `yaml:"temp_folder"`
	EnforceTargets map[string][]string `yaml:"enforce_targets"`
	Repos map[string](map[string]interface{}) `yaml:"repos"`
	Targets map[string]map[string]Target `yaml:"targets"`
}

type Param struct {
	Desc string `yaml:"desc"`
	Type string `yaml:"type"`
	Properties map[string](map[string]*Param) `yaml:"properties"`
}

type ActionFile struct {
	Version string `yaml:"version"`
	Type string `yaml:"type"`
	Params map[string](map[string]*Param) `yaml:"params"`
	Spec map[string]interface{} `yaml:"spec"`
}

type RunTimeEnv struct {
	Name string
	Version string
}

func NewRunCommand() *cobra.Command {

	var noEnforceTargets bool
	var keepTempFiles bool
	var destination string
	var gitSSHKey string

	runCmd := &cobra.Command{
		Use:   "run <namespace> <target> [gattaifile_path]",
		//Aliases: []string{"insp"},
		Short:  "Run a target",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			gattaifile_path := "GattaiFile.yaml"

			if len(args) >= 3 {
				gattaifile_path = args[2]
			}

			namespace_id := args[0]
			target_id := args[1]

			var gattaiFile GattaiFile

			yamlFile, err := ioutil.ReadFile(gattaifile_path)
			if err != nil {
				log.Fatalf("Error reading Gattai File: %v", err)
			}
			err = yaml.Unmarshal(yamlFile, &gattaiFile)
			if err != nil {
				log.Fatalf("Error parsing Gattai File: %v", err)
			}

			if noEnforceTargets == false {
				for namespace_id, target_id_list := range gattaiFile.EnforceTargets {
					if targets, ok := gattaiFile.Targets[namespace_id]; ok {
						for _, target_id := range target_id_list {
							if _, ok := targets[target_id]; !ok {
								log.Fatalf("Target from <%v> is required by enforced-target: %v", namespace_id, target_id)
							}
						}
					} else {
						log.Fatalf("Namespace is required by enforced-target: %v", namespace_id)
					}
				}
			}

			lookUpReturn := make(map[string]string)
			lookUpRepoPath := make(map[string]string)

			for key, val := range gattaiFile.Repos {
				repo, ok := val["repo"].(string)
				if ok == false {
					log.Fatalln("Please provide a repo: type")
				}
				src, ok := val["src"]
				if ok == false {
					log.Fatalln("Please provide a src: map")
				}
				switch repo {
				case "local":
					switch srcMap := src.(type){
					case map[interface{}]interface{}:
						dir, ok := srcMap["dir"].(string)
						if ok == false {
							log.Fatalln("Please provide a dir: path")
						}
						fileInfo, err := os.Stat(dir)
						if err != nil || fileInfo.IsDir() == false {
							log.Fatalln("Please provide a directory for local repo!")
						}
						lookUpRepoPath[key] = dir
					default:
						log.Fatalln("Local repo require a dir: path!")
					}
				case "git":
					switch srcMap := src.(type){
					case map[interface{}]interface{}:
						web_url, ok := srcMap["url"].(string)
						if ok == false {
						 log.Fatalln("Please provide a url: key")
						}
						_, err = url.ParseRequestURI(web_url)
						if err != nil {
							log.Fatalf("GIT repo parse request url error: %v", err)
						}
						repoDir, err := os.MkdirTemp("",key)
						if err != nil {
							log.Fatalf("Error creating repository folder: %v", err)
						}
						var ref_name plumbing.ReferenceName
						if branch, ok := srcMap["branch"].(string); ok {
							ref_name = plumbing.NewBranchReferenceName(branch)
						}
						if tag, ok := srcMap["tag"].(string); ok {
							ref_name = plumbing.NewTagReferenceName(tag)
						}
						defer os.RemoveAll(repoDir) // clean up
						_, err = git.PlainClone(repoDir, false, &git.CloneOptions{
							URL:               web_url,
							Progress: 		   os.Stdout,
							ReferenceName:	   ref_name,
						})
						if err != nil {
							log.Fatalf("Error cloning git repository: %v", err)
						}
						lookUpRepoPath[key] = repoDir
					default:
						log.Fatalln("Local repo require a dir: path!")
					}
				default:
					log.Fatalln("Repo type is not supported!")
				}
			}

			tempDir, err := os.MkdirTemp(gattaiFile.TempFolder, "gattaitmp")
			if err != nil {
				log.Fatalf("Error creating temporary folder: %v", err)
			}
			if keepTempFiles == false {
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
		},
	}

	runCmd.Flags().BoolVarP(&noEnforceTargets, "no-enforce", "n", false, "Do not enforce target")
	runCmd.Flags().BoolVarP(&keepTempFiles, "keep-temp", "k", false, "Keep temporary created files")
	runCmd.Flags().StringVarP(&destination, "destination", "d", "", "Save to filepath")
	runCmd.Flags().StringVarP(&gitSSHKey, "git-ssh-key", "g", "", "Private SSH key for git repo")

	return runCmd
}

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "gattai",
		Version: "0.0.1",
		Short: "gattai - a simple CLI to transform and inspect strings",
		Long: `gattai is a super fancy CLI (kidding)

	One can use stringer to modify or inspect strings straight from the terminal`,
		Run: func(cmd *cobra.Command, args []string) {},
	}

	rootCmd.AddCommand(NewRunCommand());

	return rootCmd
}

func val_type(item interface{}) string {
	switch i_type := item.(type) {
	case bool:
		return "bool"
	case int:
		return "int"
	case int8:
		return "int"
	case int16:
		return "int"
	case int32:
		return "int"
	case int64:
		return "int"
	case uint:
		return "int"
	case uint8:
		return "int"
	case uint16:
		return "int"
	case uint32:
		return "int"
	case uint64:
		return "int"
	case float32:
		return "float"
	case float64:
		return "float"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[interface{}]interface{}:
		return "object"
	default:
		log.Fatalf("Unsupported type: %T!\n", i_type)
	}
	return ""
}

func check_params(target Target,param_map map[string]*Param) {
	for key, val := range param_map {
		if var_item, ok := target.Vars[key]; ok {
			var_type := val_type(var_item)
			if val.Type == var_type {
				if val.Type == "object" {
					if required_params, ok := val.Properties["required"]; ok {
						check_params(target,required_params)
					}
				}
			} else {
				log.Fatalf("Invalid type for %s: %v, Expecting %v",key,var_type,val.Type)
			}
		} else {
			log.Fatalf("Missing key %s, key is required!",key)
		}
	}
}

func cmd_blk(blk interface{}) {
	switch blk_map := blk.(type) {
	case map[interface{}]interface{}:
		if command, ok := blk_map["command"].(string); ok {
			src := command
			switch args := blk_map["args"].(type) {
			case []interface {}:
				for _, elem := range args {
					src +=  " " + elem.(string)
				}
			default:
				log.Fatalf("fetch do not support type %T!\n", args)
			}

			cmd := exec.Command(command)
			_, err := cmd.Output()
			if err != nil {
				switch updated_target.RunTimeEnv {
					//case "docker":
					//	if docker, ok := rtenv_map["docker"]; ok {
					//		if app_docker, ok :=  docker[command]; ok {
					//			src = fmt.Sprintf("docker run --add-host=host.docker.internal:host-gateway --rm %s:%s%s",app_docker.Name, app_docker.Version, src)
					//		}
					//	}
					case "nix_shell":
						if nix, ok := rtenv_map["nix_shell"]; ok {
							if app_nix, ok :=  nix[command]; ok {
							src = fmt.Sprintf("nix-shell -p %s -I nixpkgs=%s --command \"%s\"", app_nix.Name, app_nix.Version, tpl_format()(src))
							}
						}
					default:
					}
			}
			//fmt.Println(src)

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
				interp.StdIO(os.Stdin, &buf, os.Stdout),
				interp.OpenHandler(open),
				interp.ExecHandler(exec),
			)
			err = runner.Run(context.TODO(), file)
			if err != nil {
				log.Fatalf("Run: %v", err)
			}
			result += buf.String()
		}
		if include_path, ok := blk_map["include"].(string); ok {
			switch vars_map := blk_map["vars"].(type) {
			case map[interface{}]interface{}:
				vars_remap := make(map[string]interface{})
				for key, val := range vars_map {
					if key_id, ok := key.(string); ok {
						vars_remap[key_id] = val
					}
				}
				new_target := Target {
					Exec: include_path,
					RunTimeEnv: updated_target.RunTimeEnv,
					Vars: vars_remap,
				}
				result += rec_cmds(new_target,repo_path, include_path, temp_folder, temp_dir)
			default:
				log.Fatalf("fetch do not support type %T!\n", vars_map)
			}
		}
	default:
		log.Fatalf("fetch do not support type %T!\n", blk_map)
	}
}

func rec_cmds(updated_target Target,repo_path string, exec_path string, temp_folder string, temp_dir string) string {

	tmpl_filepath := path.Join(repo_path,exec_path) + ".yaml"
	tmpl_filename := path.Base(tmpl_filepath)
	tmpl, err := template.New(tmpl_filename).Funcs(template.FuncMap{
		"temp_dir": tpl_temp_dir(temp_dir),
		"format": tpl_format(),
	}).ParseFiles(tmpl_filepath)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, updated_target); err != nil {
		panic(err)
	}
	var actionFile ActionFile;

	//fmt.Printf("%s\n", buf.String())
	err = yaml.Unmarshal(buf.Bytes(), &actionFile)
	if err != nil {
		log.Fatalf("Unmarshal3: %v", err)
	}

	if required_params, ok := actionFile.Params["required"]; ok {
		check_params(updated_target, required_params)
	}

	rtenv_map := make(map[string](map[string]RunTimeEnv))

	if rtEnv, ok := actionFile.Spec["runtime_env"]; ok {
		switch rt_map := rtEnv.(type){
		case map[interface{}]interface{}:
			for rt_type, rt_val := range rt_map {
				if rt_type_id, ok := rt_type.(string); ok {
					rtenv_map[rt_type_id] = make(map[string]RunTimeEnv)
					switch app_map := rt_val.(type){
					case map[interface{}]interface{}:
						for app_type, app_val := range app_map {
							if app_type_id, ok := app_type.(string); ok {
								switch cfg_map := app_val.(type){
								case map[interface{}]interface{}:
									var name, version string
									if n, ok := cfg_map["name"].(string); ok {
										name = n
									}
									if v, ok := cfg_map["version"].(string); ok {
										version = v
									}
									rtenv_map[rt_type_id][app_type_id] = RunTimeEnv{
										Name: name,
										Version: version,
									}
								default:
								}
							}
						}
					default:
					}
				}
			}
		default:
			log.Fatalf("runtime env do not support type %T!\n", rtEnv)
		}
	}

	switch test := actionFile.Spec["test"].(type){
	case []interface{}:

	default:
		log.Fatalf("fetch do not support type %T!\n", test)
	}

	var result string
	switch cmds := actionFile.Spec["cmds"].(type){
	case []interface{}:
		for _, blk := range cmds {
			switch blk_map := blk.(type) {
			case map[interface{}]interface{}:
				if command, ok := blk_map["command"].(string); ok {
					src := command
					switch args := blk_map["args"].(type) {
					case []interface {}:
						for _, elem := range args {
							src +=  " " + elem.(string)
						}
					default:
						log.Fatalf("fetch do not support type %T!\n", args)
					}

					cmd := exec.Command(command)
					_, err := cmd.Output()
					if err != nil {
						switch updated_target.RunTimeEnv {
							//case "docker":
							//	if docker, ok := rtenv_map["docker"]; ok {
							//		if app_docker, ok :=  docker[command]; ok {
							//			src = fmt.Sprintf("docker run --add-host=host.docker.internal:host-gateway --rm %s:%s%s",app_docker.Name, app_docker.Version, src)
							//		}
							//	}
							case "nix_shell":
								if nix, ok := rtenv_map["nix_shell"]; ok {
									if app_nix, ok :=  nix[command]; ok {
									src = fmt.Sprintf("nix-shell -p %s -I nixpkgs=%s --command \"%s\"", app_nix.Name, app_nix.Version, tpl_format()(src))
									}
								}
							default:
							}
					}
					//fmt.Println(src)

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
						interp.StdIO(os.Stdin, &buf, os.Stdout),
						interp.OpenHandler(open),
						interp.ExecHandler(exec),
					)
					err = runner.Run(context.TODO(), file)
					if err != nil {
						log.Fatalf("Run: %v", err)
					}
					result += buf.String()
				}
				if include_path, ok := blk_map["include"].(string); ok {
					switch vars_map := blk_map["vars"].(type) {
					case map[interface{}]interface{}:
						vars_remap := make(map[string]interface{})
						for key, val := range vars_map {
							if key_id, ok := key.(string); ok {
								vars_remap[key_id] = val
							}
						}
						new_target := Target {
							Exec: include_path,
							RunTimeEnv: updated_target.RunTimeEnv,
							Vars: vars_remap,
						}
						result += rec_cmds(new_target,repo_path, include_path, temp_folder, temp_dir)
					default:
						log.Fatalf("fetch do not support type %T!\n", vars_map)
					}
				}
			default:
				log.Fatalf("fetch do not support type %T!\n", blk_map)
			}
		}
	default:
		log.Fatalf("fetch do not support type %T!\n", cmds)
	}

	return result
}

func tpl_fetch(gattai_file GattaiFile, temp_dir string, lookUpRepoPath map[string]string, lookUpReturn map[string]string) func(Target) string {
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
				log.Fatalln("Repo prefix does not exist!")
			}

			result = strings.TrimSpace(rec_cmds(updated_target,repo_path,path.Join(tokens[1:]...),gattai_file.TempFolder,temp_dir))
			lookUpReturn[string(yamlTarget)] = result
		}
		return result
	}
}

func tpl_format() func(string) string {
	return func(content string) string {
		new_content := strings.ReplaceAll(content, "\"","\\\"")
		return strings.ReplaceAll(new_content, "\n","\\n")
	}
}

func tpl_temp_dir(temp_dir string) func(string) string {
	return func(filename string) string {
		return path.Join(temp_dir,filename)
	}
}

func main() {
	print.PrintHello()

	rootCmd := NewRootCommand()

    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
        os.Exit(1)
    }
}

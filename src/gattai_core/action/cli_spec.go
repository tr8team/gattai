package action

import (
	"os"
	"io"
	"fmt"
	"log"
	"time"
	"path"
	"bytes"
	"strings"
	"context"
	"runtime"
	"text/template"
	"gopkg.in/yaml.v2"
	//"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
	"github.com/tr8team/gattai/src/gattai_core/common"
)

const (
	NixShell string = "nix_shell"
)

const (
	CmpEqual string = "equal"
	CmpNotEqual		= "not_equal"
	CmpContain		= "contain"
	CmpNotContain	= "not_contain"
	CmpIntLessThan  = "int_less_than"
	CmpIntMoreThan  = "int_more_than"
)

type CommandLineInteraceSpec struct {
	RunTimeEnv map[string](
		map[string] struct {
			Name string `yaml:"name"`
			Version string `yaml:"version"`
		}) `yaml:"runtime_env"`
	Test struct {
		Expected struct {
			Condition string `yaml:"condition"`
			Value string `yaml:"value"`
		}
		Cmds []CmdBlock `yaml:"cmds"`
	} `yaml:"test"`
	Exec struct {
		Cmds []CmdBlock `yaml:"cmds"`
	} `yaml:"exec"`
}

type CmdBlock struct {
	Command string `yaml:"command"`
	Args [] string `yaml:"args"`
}

func RecCmds(updated_target common.Target,repo_path string, exec_path string, temp_folder string, temp_dir string) string {

	tmpl_filepath := path.Join(repo_path,exec_path) + ".yaml"
	tmpl_filename := path.Base(tmpl_filepath)
	tmpl, err := template.New(tmpl_filename).Funcs(template.FuncMap{
		"temp_dir": TplTempDir(temp_dir),
		"format": TplFormat(),
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

	CheckParams(updated_target, actionFile.Params.Required)
	var result string

	yamlSpec, err := yaml.Marshal(actionFile.Spec)
	if err != nil {
		panic(err)
	}

	switch actionFile.Type {
	case CLISpec:
		var cliSpec CommandLineInteraceSpec
		err = yaml.Unmarshal(yamlSpec, &cliSpec)
		if err != nil {
			log.Fatalf("Unmarshal4 CLISpec: %v", err)
		}

		//rtenv_map := cliSpec.RunTimeEnv

		//expected := CmdBlk(cliSpec.Test.Cmds)
		//switch cliSpec.Test.Expected.Condition {
		//case CmpEqual:
		//	if expected == cliSpec.Test.Expected.Value {
		//	}
		//case CmpNotEqual:
		//	if expected != cliSpec.Test.Expected.Value {
		//
		//	}
		//case CmpContain:
		//	if strings.Contains(expected, cliSpec.Test.Expected.Value) {
		//
		//	}
		//case CmpNotContain:
		//	if !strings.Contains(expected, cliSpec.Test.Expected.Value) {
		//
		//	}
		//case CmpIntLessThan:
		//	exp_int, err := strconv.Atoi(expected)
		//	if  err != nil {
		//	}
		//	exp_val, err := strconv.Atoi(cliSpec.Test.Expected.Value)
		//	if  err != nil {
		//	}
		//	if exp_int < exp_val {
		//
		//	}
		//case CmpIntMoreThan:
		//	exp_int, err := strconv.Atoi(expected)
		//	if  err != nil {
		//	}
		//	exp_val, err := strconv.Atoi(cliSpec.Test.Expected.Value)
		//	if  err != nil {
		//	}
		//	if exp_int > exp_val {
		//
		//	}
		//default:
		//}

		result += CmdBlk(cliSpec.Exec.Cmds)

	case WrapSpec:
		var wrapSpec WrapperInterfaceSpec
		err = yaml.Unmarshal(yamlSpec, &wrapSpec)
		if err != nil {
			log.Fatalf("Unmarshal4 wrapSpec: %v", err)
		}
		result += RecCmds(wrapSpec.Include,repo_path, wrapSpec.Include.Action, temp_folder, temp_dir)
	default:
		log.Fatalf("Action file type is not supported: %s!", actionFile.Type)
	}

	return result
}

func CmdBlk(cmds []CmdBlock) string {

	var result string

	for _, blk := range cmds {

		src := blk.Command

		for _, elem := range blk.Args {
			src +=  " " + elem
		}

		// cmd := exec.Command(blk.Command)
		// _, err := cmd.Output()
		// if err != nil {
		// 	if nix, ok := rtenv_map[NixShell]; ok {
		// 		if app_nix, ok :=  nix[blk.Command]; ok {
		// 		src = fmt.Sprintf("nix-shell -p %s -I nixpkgs=%s --command \"%s\"", app_nix.Name, app_nix.Version, tpl_format()(src))
		// 		}
		// 	}
		// }
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
			if _, err := interp.LookPathDir(hc.Dir, hc.Env, args[0]); err != nil {
				fmt.Printf("%s is not installed\n", args[0])
				return interp.NewExitStatus(1)
			}
			return interp.DefaultExecHandler(2*time.Second)(ctx, args)
		}
		var buf bytes.Buffer
		runner, _ := interp.New(
			//interp.Env(expand.ListEnviron("GLOBAL=global_value")),
			interp.StdIO(os.Stdin, &buf, os.Stdout),
			interp.OpenHandler(open),
			interp.ExecHandler(exec),
		)
		err := runner.Run(context.TODO(), file)
		if err != nil {
			log.Fatalf("Run: %v", err)
		}
		result += buf.String()
	}

	return result
}

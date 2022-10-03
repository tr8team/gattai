package action

import (
	"os"
	"io"
	"fmt"
	"log"
	"time"
	"bytes"
	"strconv"
	"strings"
	"context"
	"runtime"
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

func RunCmdBlks(cmds []CmdBlock) string {

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

func ExpectedTest(expected string, cliSpec CommandLineInteraceSpec) bool {
	result := false
	switch cliSpec.Test.Expected.Condition {
	case CmpEqual:
		result = (strings.TrimSpace(expected) == strings.TrimSpace(cliSpec.Test.Expected.Value))
		fmt.Println(result)
	case CmpNotEqual:
		result = (strings.TrimSpace(expected) != strings.TrimSpace(cliSpec.Test.Expected.Value))
	case CmpContain:
		result = strings.Contains(strings.TrimSpace(expected), strings.TrimSpace(cliSpec.Test.Expected.Value))
	case CmpNotContain:
		result =!strings.Contains(strings.TrimSpace(expected), strings.TrimSpace(cliSpec.Test.Expected.Value))
	case CmpIntLessThan:
		exp_int, err := strconv.Atoi(expected)
		if  err != nil {
			log.Fatalf("Error converting to integer: %s\n", expected)
		}
		exp_val, err := strconv.Atoi(cliSpec.Test.Expected.Value)
		if  err != nil {
			log.Fatalf("Error converting to integer: %s\n", cliSpec.Test.Expected.Value)
		}
		result = (exp_int < exp_val)
	case CmpIntMoreThan:
		exp_int, err := strconv.Atoi(expected)
		if  err != nil {
			log.Fatalf("Error converting to integer: %s\n", expected)
		}
		exp_val, err := strconv.Atoi(cliSpec.Test.Expected.Value)
		if  err != nil {
			log.Fatalf("Error converting to integer: %s\n", cliSpec.Test.Expected.Value)
		}
		result = (exp_int > exp_val)
	default:
		log.Fatalf("Condition is not supported: %s\n", cliSpec.Test.Expected.Condition)
	}
	return result
}

func ExecCLI(updated_target common.Target,actionFile ActionFile,action_args *ActionArgs) string {
	cliSpec := NewSpec[CommandLineInteraceSpec](actionFile)
	return RunCmdBlks(cliSpec.Exec.Cmds)
}

func TestCLI(updated_target common.Target,actionFile ActionFile,action_args *ActionArgs) string {
	cliSpec := NewSpec[CommandLineInteraceSpec](actionFile)
	if len(cliSpec.Test.Cmds) > 0 {
		expected := RunCmdBlks(cliSpec.Test.Cmds)
		if ExpectedTest(expected,*cliSpec) {
			log.Printf("%s: Test Passed!\n",updated_target.Action)
		} else {
			log.Fatalf("%s: Test Failed! Expecting %s, Got %s\n",updated_target.Action,cliSpec.Test.Expected.Value,expected)
		}
	} else {
		log.Printf("%s: No Test Found!\n",updated_target.Action)
	}
	return RunCmdBlks(cliSpec.Exec.Cmds)
}

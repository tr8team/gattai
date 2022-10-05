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

func ExecCommand(src string) (string, error) {
	var result bytes.Buffer
	file, err := syntax.NewParser().Parse(strings.NewReader(src), "")
	if err != nil {
		return result.String(), fmt.Errorf("ExecCommand syntaxParse error: %v",err)
	}
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
	runner, err := interp.New(
		//interp.Env(expand.ListEnviron("GLOBAL=global_value")),
		interp.StdIO(os.Stdin, &result, os.Stdout),
		interp.OpenHandler(open),
		interp.ExecHandler(exec),
	)
	if err != nil {
		return result.String(), fmt.Errorf("ExecCommand interpNew error: %v",err)
	}
	err = runner.Run(context.TODO(), file)
	if err != nil {
		return result.String(), fmt.Errorf("ExecCommand runnerRun error: %v",err)
	}
	return result.String(), nil
}

func ConstructCommand(blk CmdBlock) string {
	result := blk.Command

	for _, elem := range blk.Args {
		result +=  " " + elem
	}

	return result
}

func ExecCmdBlks(cmds []CmdBlock) (string, error) {

	var result string

	for _, blk := range cmds {

		src := ConstructCommand(blk)

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

		output, err := ExecCommand(src)
		if err != nil {
			return result, fmt.Errorf("ExecCmdBlks ExecCommand error: %v",err)
		}
		result += output
	}

	return result, nil
}

func ExpectedTest(expected string, conditon string, expected_value string) (bool,error) {
	result := false
	switch conditon {
	case CmpEqual:
		result = (strings.TrimSpace(expected) == strings.TrimSpace(expected_value))
	case CmpNotEqual:
		result = (strings.TrimSpace(expected) != strings.TrimSpace(expected_value))
	case CmpContain:
		result = strings.Contains(strings.TrimSpace(expected), strings.TrimSpace(expected_value))
	case CmpNotContain:
		result =!strings.Contains(strings.TrimSpace(expected), strings.TrimSpace(expected_value))
	case CmpIntLessThan:
		exp_int, err := strconv.Atoi(expected)
		if  err != nil {
			return result, fmt.Errorf("ExpectedTest strconvAtoi error: %s error: %v",expected, err)
		}
		exp_val, err := strconv.Atoi(expected_value)
		if  err != nil {
			return result, fmt.Errorf("ExpectedTest strconvAtoi error: %s error: %v",expected_value, err)
		}
		result = (exp_int < exp_val)
	case CmpIntMoreThan:
		exp_int, err := strconv.Atoi(expected)
		if  err != nil {
			return result, fmt.Errorf("ExpectedTest strconvAtoi error: %s error: %v",expected, err)
		}
		exp_val, err := strconv.Atoi(expected_value)
		if  err != nil {
			return result, fmt.Errorf("ExpectedTest strconvAtoi error: %s error: %v",expected_value, err)
		}
		result = (exp_int > exp_val)
	default:
		return result, fmt.Errorf("ExpectedTest condition is not supported error: %s",conditon)
	}
	return result, nil
}

func ExecCLI(updated_target common.Target,actionFile ActionFile,action_args *ActionArgs) (string,error) {
	cliSpec, err := NewSpec[CommandLineInteraceSpec](actionFile)
	if err != nil {
		return "", fmt.Errorf("ExecCLI NewSpec error: %v",err)
	}
	return ExecCmdBlks(cliSpec.Exec.Cmds)
}

func TestCLI(updated_target common.Target,actionFile ActionFile,action_args *ActionArgs) (string,error) {
	cliSpec, err := NewSpec[CommandLineInteraceSpec](actionFile)
	if err != nil {
		return "", fmt.Errorf("TestCLI NewSpec error: %v",err)
	}
	if len(cliSpec.Test.Cmds) > 0 {
		expected,err := ExecCmdBlks(cliSpec.Test.Cmds)
		if err != nil {
			return "", fmt.Errorf("TestCLI ExecCmdBlks error: %v",err)
		}
		passed, err := ExpectedTest(expected,cliSpec.Test.Expected.Condition,cliSpec.Test.Expected.Value)
		if err != nil {
			return "", fmt.Errorf("TestCLI ExpectedTest error: %v",err)
		}
		if passed {
			log.Printf("TestCLI Test Passed!:	%s\n",updated_target.Action)
		} else {
			log.Fatalf("TestCLI Test Failed!:	%s (Expecting: %s, Result: %s)\n",updated_target.Action,cliSpec.Test.Expected.Value,expected)
		}
	} else {
		log.Printf("TestCLI No Test Found!:	%s\n",updated_target.Action)
	}
	return ExecCmdBlks(cliSpec.Exec.Cmds)
}

package core_cli

import (
	"os"
	"io"
	"fmt"
	//"log"
	"time"
	"bytes"
	"strings"
	"context"
	"runtime"
	//"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
	"github.com/tr8team/gattai/src/gattai_core/core_action"
)

type CLICommand struct {
	Shell string
	EnvVars map[string]string
	CmdArray []string
}

type CLITest struct {
	Expected core_action.Comparison
	Commands []CLICommand
}

type CLIExec struct {
	Commands []CLICommand
}

func BuildCLICommand(cmd CLICommand) string {
	separator := " "
	return strings.Join(cmd.CmdArray, separator)
}

func ExecCLICommand(src string) (string, error) {
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

func RunCLICommand(cmds []CLICommand) (string, error) {

	var result string

	for _, cmd := range cmds {

		src := BuildCLICommand(cmd)

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

		output, err := ExecCLICommand(src)
		if err != nil {
			return result, fmt.Errorf("RunCLICommand ExecCLICommand error: %v",err)
		}
		result += output
	}

	return result, nil
}

func (test CLITest) RunAcion(action_name string) (string,error)  {
	result := fmt.Sprintf("%s No Test Found!\n",action_name)
	if len(test.Commands) > 0 {
		expected,err := RunCLICommand(test.Commands)
		if err != nil {
			return result, fmt.Errorf("%s TestAction RunCLICommand error: %v",action_name,err)
		}
		passed, err := core_action.ExpectedTest(expected,test.Expected.Condition,test.Expected.Value)
		if err != nil {
			return result, fmt.Errorf("%s TestAction ExpectedTest error: %v",action_name,err)
		}
		if passed {
			result = fmt.Sprintf("%s Test Passed!\n",action_name)
		} else {
			return result, fmt.Errorf("%s TestAction Test Failed! (Expecting: %s, Result: %s)\n",action_name,test.Expected.Value,expected)
		}
	}
	return result, nil
}

func (exec CLIExec) RunAcion(action_name string) (string,error)  {
	return RunCLICommand(exec.Commands)
}

package yaml_format

import (
	"os"
	"io"
	"fmt"
	//"log"
	"time"
	"bytes"
	//"strconv"
	"strings"
	"context"
	"runtime"
	//"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
	"github.com/tr8team/gattai/src/gattai_core/core_action"
	"github.com/tr8team/gattai/src/gattai_core/core_cli"
)

const (
	NixShell string = "nix_shell"
)



type CommandLineInteraceSpec struct {
	// RunTimeEnv map[string](
	// 	map[string] struct {
	// 		Name string `yaml:"name"`
	// 		Version string `yaml:"version"`
	// 	}) `yaml:"runtime_env"`
	Test TestCmd `yaml:"test"`
	Exec struct {
		Cmds []CmdBlock `yaml:"cmds"`
	} `yaml:"exec"`
}

type TestCmd struct {
	Expected struct {
		Condition string `yaml:"condition"`
		Value string `yaml:"value"`
	}
	Cmds []CmdBlock `yaml:"cmds"`
}

type CmdBlock struct {
	Command string `yaml:"command"`
	Args [] string `yaml:"args"`
}

func (blk CmdBlock) GetArray()[]string{
	return append(
		[]string{blk.Command},
		blk.Args...
	)
}

func ConvertToCLICommand(shell string, envVars map[string]string, arr []CmdBlock) []core_cli.CLICommand {
	result := make([]core_cli.CLICommand, len(arr))
	for i, blk := range arr {
		result[i] = core_cli.CLICommand {
			Shell: shell,
			EnvVars: envVars,
			CmdArray: blk.GetArray(),
		}
	}
	return result
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

func (cliSpec CommandLineInteraceSpec) GenerateAction(action_name string, action_args ActionArgs) (*core_action.Action,error)  {
	return &core_action.Action {
		Name: action_name,
		Test: &core_cli.CLITest {
			Expected: core_action.Comparison {
				Condition: cliSpec.Test.Expected.Condition,
				Value: cliSpec.Test.Expected.Value,
			},
			Commands: ConvertToCLICommand("",make(map[string]string),cliSpec.Test.Cmds),
		},
		Exec: core_cli.CLIExec {
			Commands: ConvertToCLICommand("",make(map[string]string),cliSpec.Exec.Cmds),
		},
	}, nil
}

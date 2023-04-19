package yaml_format

import (
	"fmt"
	"path"
	"github.com/tr8team/gattai/src/gattai_core/cli"
)

type DerivedInterfaceSpec struct {
	Repo Repo `yaml:"repo"`
	InheritExec Target `yaml:"inherit_exec"`
	OverrideTest TestCmd `yaml:"override_test"`
}

func (diSpec * DerivedInterfaceSpec) Derived(action_name string,action_args ActionArgs) (ActionSpecInterface,error) {
	repopath := action_args.RepoPath
	if len(diSpec.Repo.Src) > 0 {
		output, err := GetRepoPath(action_args.TempDir,diSpec.Repo.Src,diSpec.Repo,action_args.RepoPath)
		if err != nil {
			return nil, fmt.Errorf("%s error: %v",action_name,err)
		}
		repopath = output
	}
	tmpl_filepath := path.Join(repopath,diSpec.InheritExec.Action) + ".yaml"
	return RunAction(diSpec.InheritExec,tmpl_filepath,action_args)
}

func (diSpec DerivedInterfaceSpec) GenerateTestAction(action_name string, action_args ActionArgs) (*cli.CLIAction,error)  {
	if len(diSpec.OverrideTest.Cmds) > 0 {
		return &cli.CLIAction{
			Expected: cli.Comparison {
				Condition: diSpec.OverrideTest.Expected.Condition,
				Value: diSpec.OverrideTest.Expected.Value,
			},
			Exec: func(arr []CmdBlock) []cli.CLICommand {
				result := make([]cli.CLICommand, len(arr))
				for i, blk := range arr {
					result[i] = cli.CLICommand {
						Shell: "",
						EnvVars: make(map[string]string),
						CmdArray: blk.GetArray(),
					}
				}
				return result
			}(diSpec.OverrideTest.Cmds),
		}, nil
	}
	actSpec, err := diSpec.Derived(action_name,action_args)
	if err != nil {
		return nil, fmt.Errorf("%s GenerateTestAction error: %v",action_name,err)
	}
	return actSpec.GenerateTestAction(action_name,action_args)
}

func (diSpec DerivedInterfaceSpec) GenerateExecAction(action_name string, action_args ActionArgs) (*cli.CLIAction,error)  {
	actSpec, err := diSpec.Derived(action_name,action_args)
	if err != nil {
		return nil, fmt.Errorf("%s GenerateExecAction error: %v",action_name,err)
	}
	return actSpec.GenerateExecAction(action_name,action_args)
}

// func (diSpec DerivedInterfaceSpec) TestAction(action_name string, action_args ActionArgs) (string,error)  {

// 	if len(diSpec.OverrideTest.Cmds) > 0 {
// 		expected, err := ExecCmdBlks(diSpec.OverrideTest.Cmds)
// 		if err != nil {
// 			return "", fmt.Errorf("%s ExecCmdBlks error: %v",action_name,err)
// 		}
// 		passed, err := ExpectedTest(expected,diSpec.OverrideTest.Expected.Condition,diSpec.OverrideTest.Expected.Value)
// 		if err != nil {
// 			return "", fmt.Errorf("%s ExpectedTest error: %v",action_name,err)
// 		}
// 		if passed {
// 			return fmt.Sprintf("%s Test Passed!\n",action_name), nil
// 		} else {
// 			return "", fmt.Errorf("%s Test Failed! (Expecting: %s, Result: %s)\n",action_name,diSpec.OverrideTest.Expected.Value,expected)
// 		}
// 	}
// 	actSpec, err := diSpec.Derived(action_name,action_args)
// 	if err != nil {
// 		return "", fmt.Errorf("%s TestAction error: %v",action_name,err)
// 	}
// 	return actSpec.TestAction(action_name,action_args)
// }

// func (diSpec DerivedInterfaceSpec) ExecAction(action_name string, action_args ActionArgs) (string,error) {
// 	actSpec, err := diSpec.Derived(action_name,action_args)
// 	if err != nil {
// 		return "", fmt.Errorf("%s ExecAction error: %v",action_name,err)
// 	}
// 	return actSpec.ExecAction(action_name,action_args)
// }

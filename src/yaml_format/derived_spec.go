package yaml_format

import (
	"fmt"
	"path"
	"github.com/tr8team/gattai/src/gattai_core/core_action"
	"github.com/tr8team/gattai/src/gattai_core/core_cli"
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

func (diSpec DerivedInterfaceSpec) GenerateTestAction(action_name string, action_args ActionArgs) (core_action.ActionInterface,error)  {
	if len(diSpec.OverrideTest.Cmds) > 0 {
		return &core_cli.CLIAction{
			Expected: core_action.Comparison {
				Condition: diSpec.OverrideTest.Expected.Condition,
				Value: diSpec.OverrideTest.Expected.Value,
			},
			Exec: func(arr []CmdBlock) []core_cli.CLICommand {
				result := make([]core_cli.CLICommand, len(arr))
				for i, blk := range arr {
					result[i] = core_cli.CLICommand {
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

func (diSpec DerivedInterfaceSpec) GenerateExecAction(action_name string, action_args ActionArgs) (core_action.ActionInterface,error)  {
	actSpec, err := diSpec.Derived(action_name,action_args)
	if err != nil {
		return nil, fmt.Errorf("%s GenerateExecAction error: %v",action_name,err)
	}
	return actSpec.GenerateExecAction(action_name,action_args)
}

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
	return GenerateActionSpec(diSpec.InheritExec,tmpl_filepath,action_args)
}

func (diSpec DerivedInterfaceSpec) GenerateAction(action_name string, action_args ActionArgs) (*core_action.Action,error)  {
	actSpec, err := diSpec.Derived(action_name,action_args)
	if err != nil {
		return nil, fmt.Errorf("%s GenerateAction Derived error: %v",action_name,err)
	}
	action,err := actSpec.GenerateAction(action_name,action_args)
	if err != nil {
		return nil, fmt.Errorf("%s GenerateAction GenerateAction error: %v",action_name,err)
	}
	if len(diSpec.OverrideTest.Cmds) > 0 {
		return &core_action.Action{
			Name: action_name,
			Test: core_cli.CLITest {
				Expected: core_action.Comparison {
					Condition: diSpec.OverrideTest.Expected.Condition,
					Value: diSpec.OverrideTest.Expected.Value,
				},
				Commands: ConvertToCLICommand("",make(map[string]string),diSpec.OverrideTest.Cmds),
			},
			Exec: action.Exec,
		}, nil
	}
	return action, nil
}

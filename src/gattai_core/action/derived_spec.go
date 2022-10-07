package action

import (
	"fmt"
	"path"
	"github.com/tr8team/gattai/src/gattai_core/common"
)

type DerivedInterfaceSpec struct {
	Repo common.Repo `yaml:"repo"`
	Include common.Target `yaml:"include"`
}

func RedirectDerived(updated_target common.Target,actionFile ActionFile,action_args *ActionArgs) (string,error) {
	derivedSpec,err := NewSpec[DerivedInterfaceSpec](actionFile)
	if err != nil {
		return "", fmt.Errorf("RedirectDerived error: %v",err)
	}
	repopath := action_args.RepoPath
	if len(derivedSpec.Repo.Src) > 0 {
		output, err := common.GetRepoPath(action_args.TempDir,derivedSpec.Repo.Src,derivedSpec.Repo)
		if err != nil {
			return "", fmt.Errorf("RedirectDerived error: %v",err)
		}
		repopath = output
	}
	tmpl_filepath := path.Join(repopath,derivedSpec.Include.Action) + ".yaml"
	return RunAction(derivedSpec.Include,tmpl_filepath,action_args)
}

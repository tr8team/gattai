package action

import (
	"fmt"
	"github.com/tr8team/gattai/src/gattai_core/common"
)

type DerivedInterfaceSpec struct {
	Include common.Target `yaml:"include"`
}

func RedirectDerived(updated_target common.Target,actionFile ActionFile,action_args *ActionArgs) (string,error) {
	derivedSpec,err := NewSpec[DerivedInterfaceSpec](actionFile)
	if err != nil {
		return "", fmt.Errorf("RedirectDerived NewSpec error: %v",err)
	}
	return RunAction(derivedSpec.Include,derivedSpec.Include.Action,action_args)
}

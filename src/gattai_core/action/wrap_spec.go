package action

import (
	"fmt"
	"github.com/tr8team/gattai/src/gattai_core/common"
)

type WrapperInterfaceSpec struct {
	Include common.Target `yaml:"include"`
}

func RedirectWrap(updated_target common.Target,actionFile ActionFile,action_args *ActionArgs) (string,error) {
	wrapSpec,err := NewSpec[WrapperInterfaceSpec](actionFile)
	if err != nil {
		return "", fmt.Errorf("RedirectWrap NewSpec error: %v",err)
	}
	return RunAction(wrapSpec.Include,wrapSpec.Include.Action,action_args)
}

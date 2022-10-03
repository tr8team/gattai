package action

import (
"github.com/tr8team/gattai/src/gattai_core/common"
)

type WrapperInterfaceSpec struct {
	Include common.Target `yaml:"include"`
}

func RedirectWrap(updated_target common.Target,actionFile ActionFile,action_args *ActionArgs) string {
	wrapSpec := NewSpec[WrapperInterfaceSpec](actionFile)
	return RunAction(wrapSpec.Include,wrapSpec.Include.Action,action_args)
}

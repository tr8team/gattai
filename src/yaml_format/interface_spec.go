package yaml_format

import (
	"github.com/tr8team/gattai/src/gattai_core/core_action"
)

type ActionSpecInterface interface {
	GenerateTestAction(string,ActionArgs) (core_action.ActionInterface,error)
	GenerateExecAction(string,ActionArgs) (core_action.ActionInterface,error)
}

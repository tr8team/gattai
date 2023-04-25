package yaml_format

import (
	"github.com/tr8team/gattai/src/gattai_core/core_action"
)

type ActionSpecInterface interface {
	GenerateAction(string,ActionArgs) (*core_action.Action,error)
}

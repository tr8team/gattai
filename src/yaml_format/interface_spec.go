package yaml_format

import (
	"github.com/tr8team/gattai/src/gattai_core/core_engine"
)

type ActionSpecInterface interface {
	GenerateAction(string,ActionArgs) (*core_engine.Action,error)
}

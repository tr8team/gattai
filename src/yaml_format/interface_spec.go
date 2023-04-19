package yaml_format

import (
	"github.com/tr8team/gattai/src/gattai_core/cli"
)

type ActionSpecInterface interface {
	GenerateTestAction(string,ActionArgs) (*cli.CLIAction,error)
	GenerateExecAction(string,ActionArgs) (*cli.CLIAction,error)
}

package command

import (
	"log"
	"github.com/spf13/cobra"
	"github.com/tr8team/gattai/src/gattai_core/core_cli"
	"github.com/tr8team/gattai/src/gattai_core/core_action"
	"github.com/tr8team/gattai/src/gattai_core/core_engine"
)

func TestFetchFunc(cmdTmpArray[]string, keyMap map[string]string) core_engine.FetchFunc {
	return func(targetKey string, engine *core_engine.Engine) (*core_action.Action,error) {
		result := make([]string, len(cmdTmpArray))
		for i, token := range cmdTmpArray {
			if value, ok := keyMap[token]; ok {
				result[i] = engine.Fetch(value)
			} else {
				result[i] = token
			}
		}
		return &core_action.Action{
			Exec: core_cli.CLIExec {
				Commands: []core_cli.CLICommand {
					{
						Shell: "",
						EnvVars: make(map[string]string),
						CmdArray: result,
					},
				},
			},
		}, nil
	}
}

func NewTestCommand() *cobra.Command {

	testCmd := &cobra.Command{
		Use:   "test",
		Short:  "Test manually setup call",
		//Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			engine := core_engine.MakeEngine(core_action.RunCmdAction)

			// E   D
			// \ /
			//  C   B
			//   \ /
			//    A

			engine.Store("test-A",TestFetchFunc(
				[]string{"expr","<B>","+","<C>"},
				map[string]string{
					"<B>": "test-B",
					"<C>": "test-C",
				},
			))
			engine.Store("test-B",TestFetchFunc(
				[]string{"expr 1 + 1"},
				make(map[string]string),
			))
			engine.Store("test-C",TestFetchFunc(
				[]string{"expr","<D>","+","<E>"},
				map[string]string{
					"<D>": "test-D",
					"<E>": "test-E",
				},
			))
			engine.Store("test-D",TestFetchFunc(
				[]string{"expr 2 + 2"},
				make(map[string]string),
			))
			engine.Store("test-E",TestFetchFunc(
				[]string{"expr 3 + 3"},
				make(map[string]string),
			))

			log.Println(engine.Fetch("test-A"))
		},
	}

	return testCmd
}

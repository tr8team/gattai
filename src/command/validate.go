package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tr8team/gattai/src/gattai_core"
)

func NewValidateCommand() *cobra.Command {

	validCmd := &cobra.Command{
		Use:   "validate <namespace> <target> [gattaifile_path]",
		//Aliases: []string{"insp"},
		Short:  "Validate a target",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			gattaifile_path := "GattaiFile.yaml"

			if len(args) >= 3 {
				gattaifile_path = args[2]
			}

			namespace_id := args[0]
			target_id := args[1]

			gattaiFile := core.NewGattaiFile(gattaifile_path)

			gattaiFile.CheckEnforceTargets()

			lookUpRepoPath := gattaiFile.BuildRepoMap()

			tempDir := gattaiFile.CreateTempDir(false)

			lookUpReturn := make(map[string]string)
			switch namespace_id {
			case core.AllNamespaces:
				switch  target_id {
				case core.AllTargets:
					// all namespaces and all targets
					for _, targets := range gattaiFile.Targets {
						for _, target := range targets {
							result := core.TplFetch(*gattaiFile,tempDir,lookUpRepoPath,lookUpReturn)(target)
							fmt.Println(result)
						}
					}
				default:
					// all namespaces and a single target
					for _, targets := range gattaiFile.Targets {
						if target, ok := targets[target_id]; ok {
							result := core.TplFetch(*gattaiFile,tempDir,lookUpRepoPath,lookUpReturn)(target)
							fmt.Println(result)
						}
					}
				}
			default:
				if targets , ok := gattaiFile.Targets[namespace_id]; ok {
					switch  target_id {
					case core.AllTargets:
						// a single namespace and all targets
						for _, target := range targets {
							result := core.TplFetch(*gattaiFile,tempDir,lookUpRepoPath,lookUpReturn)(target)
							fmt.Println(result)
						}
					default:
						// a single namespace and a single target
						if target, ok := targets[target_id]; ok {
							result := core.TplFetch(*gattaiFile,tempDir,lookUpRepoPath,lookUpReturn)(target)
							fmt.Println(result)
						}
					}
				}
			}
		},
	}

	return validCmd
}

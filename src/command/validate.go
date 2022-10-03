package command

import (
	"os"
	"fmt"
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
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

			var gattaiFile core.GattaiFile

			yamlFile, err := ioutil.ReadFile(gattaifile_path)
			if err != nil {
				log.Fatalf("Error reading Gattai File: %v", err)
			}
			err = yaml.Unmarshal(yamlFile, &gattaiFile)
			if err != nil {
				log.Fatalf("Error parsing Gattai File: %v", err)
			}

			//if noEnforceTargets == false {
				for namespace_id, target_id_list := range gattaiFile.EnforceTargets {
					if targets, ok := gattaiFile.Targets[namespace_id]; ok {
						for _, target_id := range target_id_list {
							if _, ok := targets[target_id]; !ok {
								log.Fatalf("Target from <%v> is required by enforced-target: %v", namespace_id, target_id)
							}
						}
					} else {
						log.Fatalf("Namespace is required by enforced-target: %v", namespace_id)
					}
				}
			//}

			lookUpReturn := make(map[string]string)
			lookUpRepoPath := core.BuildRepoMap(gattaiFile)

			tempDir, err := os.MkdirTemp(gattaiFile.TempFolder, core.GattaiTmpFolder)
			if err != nil {
				log.Fatalf("Error creating temporary folder: %v", err)
			}
			//if keepTempFiles == false {
				fmt.Println("Clean up temp files!")
				defer os.RemoveAll(tempDir) // clean up
			//}

			switch namespace_id {
			case "*":
				switch  target_id {
				case "*":
					// all namespaces and all targets
					for _, targets := range gattaiFile.Targets {
						for _, target := range targets {
							result := core.TplFetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn)(target)
							fmt.Println(result)
						}
					}
				default:
					// all namespaces and a single target
					for _, targets := range gattaiFile.Targets {
						if target, ok := targets[target_id]; ok {
							result := core.TplFetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn)(target)
							fmt.Println(result)
						}
					}
				}
			default:
				if targets , ok := gattaiFile.Targets[namespace_id]; ok {
					switch  target_id {
					case "*":
						// a single namespace and all targets
						for _, target := range targets {
							result := core.TplFetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn)(target)
							fmt.Println(result)
						}
					default:
						// a single namespace and a single target
						if target, ok := targets[target_id]; ok {
							result := core.TplFetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn)(target)
							fmt.Println(result)
						}
					}
				}
			}
		},
	}

	return validCmd
}

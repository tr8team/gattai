package command

import (
	"os"
	"fmt"
	"log"
	"github.com/spf13/cobra"
	"github.com/tr8team/gattai/src/yaml_format"
	"github.com/tr8team/gattai/src/gattai_core/core_action"
)

func ValidateCmdAction(action *core_action.Action) (string, error){
	result, err := action.Test.RunAction(action.Name)
	if err != nil {
		return "", fmt.Errorf("ValidateCmdAction error: %v",err)
	}
	log.Println(result)
	return action.Exec.RunAction(action.Name)
}

func NewValidateCommand() *cobra.Command {

	validCmd := &cobra.Command{
		Use:   "validate <namespace> <target> [gattaifile_path|gattaifile_folder]",
		Aliases: []string{"valid"},
		Short:  "Validate a target",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			gattaifile_path := yaml_format.GattaiFileDefault

			if len(args) >= 3 {
				output, err := GetGattaiFilePath(args[2],gattaifile_path)
				if err != nil {
					log.Fatalf("Error Gattai file path: %v", err)
				}
				gattaifile_path = output
			}

			namespace_id := args[0]
			target_id := args[1]

			gattaiFile,err := yaml_format.NewGattaiFile(gattaifile_path)
			if err != nil {
				log.Fatalf("Error parsing Gattai file: %v", err)
			}

			err = gattaiFile.CheckVersion()
			if err != nil {
				log.Fatalf("Gattai version error: %v!\n", err)
			}

			err = gattaiFile.CheckEnforceTargets()
			if err != nil {
				log.Fatalln(err)
			}

			tempDir, err := gattaiFile.CreateTempDir(yaml_format.GattaiTmpFolder)
			if err != nil {
				log.Fatalf("Error creating temporary folder: %v", err)
			}
			log.Println("Clean up temp files!")
			defer os.RemoveAll(tempDir) // clean up


			gattaiFile.LookupTargets(namespace_id, target_id, tempDir,ValidateCmdAction)
		},
	}

	return validCmd
}

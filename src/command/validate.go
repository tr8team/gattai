package command

import (
	"os"
	"fmt"
	"log"
	"github.com/spf13/cobra"
	"github.com/tr8team/gattai/src/gattai_core"
	"github.com/tr8team/gattai/src/gattai_core/yaml_action"
)

func ValidateCmdAction(actSpec yaml_action.ActionSpecInterface, actArgs yaml_action.ActionArgs, actName string) (string, error){
	cliTestAct, err := actSpec.GenerateTestAction(actName,actArgs)
	if err != nil {
		return "", fmt.Errorf("ValidateCmdAction error: %v",err)
	}
	result, err := cliTestAct.TestAction(actName)
	if err != nil {
		return "", fmt.Errorf("ValidateCmdAction error: %v",err)
	}
	log.Println(result)
	cliExecAct, err := actSpec.GenerateExecAction(actName,actArgs)
	if err != nil {
		return "", fmt.Errorf("ValidateCmdAction error: %v",err)
	}
	return cliExecAct.ExecAction(actName)
}

func NewValidateCommand() *cobra.Command {

	validCmd := &cobra.Command{
		Use:   "validate <namespace> <target> [gattaifile_path|gattaifile_folder]",
		Aliases: []string{"valid"},
		Short:  "Validate a target",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			gattaifile_path := core.GattaiFileDefault

			if len(args) >= 3 {
				output, err := GetGattaiFilePath(args[2],gattaifile_path)
				if err != nil {
					log.Fatalf("Error Gattai file path: %v", err)
				}
				gattaifile_path = output
			}

			namespace_id := args[0]
			target_id := args[1]

			gattaiFile,err := core.NewGattaiFile(gattaifile_path)
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

			tempDir, err := gattaiFile.CreateTempDir(core.GattaiTmpFolder)
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

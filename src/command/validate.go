package command

import (
	"os"
	//"fmt"
	"log"
	"github.com/spf13/cobra"
	"github.com/tr8team/gattai/src/gattai_core"
	"github.com/tr8team/gattai/src/gattai_core/action"
)

func NewValidateCommand() *cobra.Command {

	validCmd := &cobra.Command{
		Use:   "validate <namespace> <target> [gattaifile_path]",
		//Aliases: []string{"insp"},
		Short:  "Validate a target",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			gattaifile_path := core.GattaiFileDefault

			if len(args) >= 3 {
				gattaifile_path = args[2]
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


			gattaiFile.LookupTargets(namespace_id, target_id, tempDir,map[string]action.ActionFunc{
				action.ActionVerKey(action.CLISpec, action.Version1): action.TestCLI,
				action.ActionVerKey(action.WrapSpec, action.Version1): action.RedirectWrap,
			})
		},
	}

	return validCmd
}

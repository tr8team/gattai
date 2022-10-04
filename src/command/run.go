package command

import (
	"os"
	"fmt"
	"log"
	"github.com/spf13/cobra"
	"github.com/tr8team/gattai/src/gattai_core"
	"github.com/tr8team/gattai/src/gattai_core/action"
)

func NewRunCommand() *cobra.Command {

	var enforceTargets bool
	var keepTempFiles bool

	runCmd := &cobra.Command{
		Use:   "run <namespace> <target> [gattaifile_path]",
		//Aliases: []string{"insp"},
		Short:  "Run a target",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			gattaifile_path := "GattaiFile.yaml"

			if len(args) >= 3 {
				gattaifile_path = args[2]
			}

			namespace_id := args[0]
			target_id := args[1]

			gattaiFile := core.NewGattaiFile(gattaifile_path)

			if gattaiFile.Version != core.Version1 {
				log.Fatalf("Gattai version not supported: %T=v!\n", gattaiFile.Version)
			}

			if enforceTargets {
				enforced_list := gattaiFile.CheckEnforceTargets()
				if len(enforced_list) > 0 {
					log.Fatalln(enforced_list)
				}
			}

			tempDir := gattaiFile.CreateTempDir()
			if keepTempFiles == false {
				fmt.Println("Clean up temp files!")
				defer os.RemoveAll(tempDir) // clean up
			}

			result := gattaiFile.LookupTargets(namespace_id, target_id, tempDir,map[string]action.ActionFunc{
				action.ActionVerKey(action.CLISpec, action.Version1): action.ExecCLI,
				action.ActionVerKey(action.WrapSpec, action.Version1): action.RedirectWrap,
			})

			fmt.Println(result)
		},
	}

	runCmd.Flags().BoolVarP(&enforceTargets, "enforce", "e", false, "Run enforce target")
	runCmd.Flags().BoolVarP(&keepTempFiles, "keep-temp", "k", false, "Keep temporary created files")

	return runCmd
}
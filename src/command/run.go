package command

import (
	"os"
	"fmt"
	"log"
	"path"
	"github.com/spf13/cobra"
	"github.com/tr8team/gattai/src/yaml_format"
	"github.com/tr8team/gattai/src/gattai_core/core_action"
)

func NewRunCommand() *cobra.Command {

	var enforceTargets bool
	var keepTempFiles bool

	runCmd := &cobra.Command{
		Use:   "run <namespace> <target> [gattaifile_path|gattaifile_folder]",
		//Aliases: []string{"insp"},
		Short:  "Run a target",
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

			fileInfo, err := os.Stat(gattaifile_path)
			if err != nil {
				log.Fatalf("Error invalid Gattai file: %v", err)
			}

			if fileInfo.IsDir() {
				gattaifile_path = path.Join(gattaifile_path,yaml_format.GattaiFileDefault)
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

			if enforceTargets {
				err := gattaiFile.CheckEnforceTargets()
				if err != nil {
					log.Fatalln(err)
				}
			}

			tempDir, err := gattaiFile.CreateTempDir(yaml_format.GattaiTmpFolder)
			if err != nil {
				log.Fatalf("Error creating temporary folder: %v", err)
			}
			if keepTempFiles == false {
				log.Printf("Clean up temp files in: %s", tempDir)
				defer os.RemoveAll(tempDir) // clean up
			}

			result,err := gattaiFile.LookupTargets(namespace_id, target_id, tempDir,core_action.RunCmdAction)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(result)
		},
	}

	runCmd.Flags().BoolVarP(&enforceTargets, "enforce", "e", false, "Run enforce target")
	runCmd.Flags().BoolVarP(&keepTempFiles, "keep-temp", "k", false, "Keep temporary created files")

	return runCmd
}

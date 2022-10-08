package command

import (
	"os"
	//"fmt"
	"log"
	"path"
	"bytes"
	"text/template"
	"gopkg.in/yaml.v2"
	"github.com/spf13/cobra"
	"github.com/tr8team/gattai/src/gattai_core/common"
	"github.com/tr8team/gattai/src/gattai_core/action"
)

func ReadActionFilesDir(input_path string) {
	fileInfo, err := os.Stat(input_path)
	if err != nil {
		log.Fatalf("Error invalid path: %v", err)
	}

	if fileInfo.IsDir() {
		items, err := os.ReadDir(input_path)
		if err != nil {
			log.Fatalf("Error cannot read dir: %v", err)
		}
		for _, item := range items {
			entry := item.Name()
			ReadActionFilesDir(path.Join(input_path,entry))
		}
	} else {
		ReadSingleActionFile(input_path)
	}
}

func ReadSingleActionFile(file_path string) {
	tmpl_filename := path.Base(file_path)
	tmpl, err := template.New(tmpl_filename).Funcs(template.FuncMap{
		"temp_dir": action.TplTempDir(""),
		"format": action.TplFormat(),
	}).ParseFiles(file_path)
	if err != nil {
		log.Fatalf("Error template error: %v",err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, common.Target{})
	if err != nil {
		log.Fatalf("Error Execute error: %v",err)
	}
	var actionFile action.ActionFile

	err = yaml.Unmarshal(buf.Bytes(), &actionFile)
	if err != nil {
		log.Fatalf("Error Unmarshal: %v", err)
	}
	err = actionFile.CheckVersion()
	if err != nil {
		log.Fatalf("Error CheckVersion: %v", err)
	}

	output, err := actionFile.GenerateTargetFromParamsInYaml()
	if err != nil {
		log.Fatalf("Error GenerateTargetFromParamsInYaml: %v", err)
	}
	log.Println(output)
}

func NewDocumentCommand() *cobra.Command {

	var recursive bool

	docCmd := &cobra.Command{
		Use:   "document [actionfile_path|actionfile_folder]",
		Aliases: []string{"doc"},
		Short:  "Document an action",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ReadActionFilesDir(args[0])
		},
	}

	docCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively generate documents")

	return docCmd
}

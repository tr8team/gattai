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

func ReadActionFilesDir(root_path string, item_name string) map[string]string {
	input_path := path.Join(root_path,item_name)
	fileInfo, err := os.Stat(input_path)
	if err != nil {
		log.Fatalf("Error invalid path: %v", err)
	}

	result := make(map[string]string)

	if fileInfo.IsDir() {
		items, err := os.ReadDir(input_path)
		if err != nil {
			log.Fatalf("Error cannot read dir: %v", err)
		}
		for _, item := range items {
			copiedmap := ReadActionFilesDir(input_path,item.Name())
			for key, val  := range copiedmap {
				result[key] = val
		    }
		}
	} else {
		if path.Ext(input_path) == ".yaml" {
			result[input_path] = ReadSingleActionFile(root_path,item_name)
		}
	}
	return result
}

func ReadSingleActionFile(root_path string, filename string) string {
	file_path := path.Join(root_path,filename)
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

	result, err := actionFile.GenerateTargetFromParamsInYaml(file_path)
	if err != nil {
		log.Fatalf("Error GenerateTargetFromParamsInYaml: %v", err)
	}
	return result
}

func NewDocumentCommand() *cobra.Command {

	var recursive bool

	docCmd := &cobra.Command{
		Use:   "document [actionfile_path|actionfile_folder]",
		Aliases: []string{"doc"},
		Short:  "Document an action",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			finalmap := ReadActionFilesDir("",args[0])
			for key, val  := range finalmap {
				log.Printf("%s : %s\n",key,val)
		    }
		},
	}

	docCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively generate documents")

	return docCmd
}

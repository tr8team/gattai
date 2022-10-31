package command

import (
	"os"
	"fmt"
	"log"
	"path"
	"bytes"
	"text/template"
	"gopkg.in/yaml.v2"
	"github.com/spf13/cobra"
	"github.com/tr8team/gattai/src/gattai_core/common"
	"github.com/tr8team/gattai/src/gattai_core/action"
)

type ReadMeDoc struct {
	Entries map[string]ReadMeEntry `yaml:"Entries"`
}

type ReadMeEntry struct {
	Fields []action.ParamField `yaml:"fields"`
	YamlTarget string `yaml:"yamlTarget"`
}

func (doc ReadMeDoc) Print() (string, error) {
	log.Printf("ReadMeDoc: %v\n",doc)
	content := `
		<table>
		<tr>
		<td> File </td> <td> Fields </td><td>Description</td>
		</tr>
		{{- range $key, $val := .Entries }}
		{{- range $index, $elem := $val.fields }}
		{{- if equal 0 $index }}
		<tr>
		<td rowspan="{{ len $val.fields }}">
		<b>{{ $key }}</b>
		\n
		{{ $val.yamlTarget }}
		\n
		</td>
		{{- else }}
		<tr>
		{{- end }}
		<td>{{ $elem.name }}<br/>{{ $elem.attribute }}</td>
		<td>{{ $elem.desc }}</td>
		</tr>
		{{- end }}
		{{- end }}
		</table>
	`
	tmpl, err := template.New("").Parse(content)
	var buf bytes.Buffer
	err = tmpl.Execute(&buf,doc)
	if err != nil {
		return "", fmt.Errorf("ReadMeDoc Print error: %v",err)
	}
	return buf.String(), nil
}

func ReadActionFilesDir(root_path string, item_name string) map[string]ReadMeEntry {
	input_path := path.Join(root_path,item_name)
	fileInfo, err := os.Stat(input_path)
	if err != nil {
		log.Fatalf("Error invalid path: %v", err)
	}

	result := make(map[string]ReadMeEntry)

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

func ReadSingleActionFile(root_path string, filename string) ReadMeEntry {
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

	yamlTarget, err := actionFile.GenerateTargetFromParamsInYaml(file_path)
	if err != nil {
		log.Fatalf("Error GenerateTargetFromParamsInYaml: %v", err)
	}

	paramField := actionFile.GenerateParamFields()

	return ReadMeEntry {
		Fields : paramField,
		YamlTarget : yamlTarget,
	}
}

func NewDocumentCommand() *cobra.Command {

	var recursive bool

	docCmd := &cobra.Command{
		Use:   "document [actionfile_path|actionfile_folder]",
		Aliases: []string{"doc"},
		Short:  "Document an action",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			finalmap := ReadMeDoc{ Entries: ReadActionFilesDir("",args[0]) }
			output, err := finalmap.Print()
			if err != nil {
				log.Fatalf("Error NewDocumentCommand: %v", err)
			}
			log.Println(output)
		},
	}

	docCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively generate documents")

	return docCmd
}

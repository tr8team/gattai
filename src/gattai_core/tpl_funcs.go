package core

import (
	//"os"
	//"io"
	//"fmt"
	"log"
	//"time"
	"path"
	"bytes"
	"strings"
	//"context"
	//"runtime"
	"text/template"
	"gopkg.in/yaml.v2"
	"github.com/tr8team/gattai/src/gattai_core/action"
	"github.com/tr8team/gattai/src/gattai_core/common"
)

func TplFetch(gattai_file GattaiFile, temp_dir string, lookUpRepoPath map[string]string, lookUpReturn map[string]string) func(common.Target) string {
	return func(target common.Target) string {
		// get target generated key
		yamlTarget, err := yaml.Marshal(target)
		if err != nil {
			panic(err)
		}
		var buf bytes.Buffer
		// check if result for target already exist
		result, ok := lookUpReturn[string(yamlTarget)]
		if !ok {
			// if not, parse target to see if target have dependency
			tmpl, err := template.New("").Funcs(template.FuncMap{
				"fetch": TplFetch(gattai_file,temp_dir,lookUpRepoPath,lookUpReturn),
			}).Parse(string(yamlTarget))
			if err != nil {
				panic(err)
			}
			buf.Reset()
			if err := tmpl.Execute(&buf, gattai_file); err != nil {
				panic(err)
			}
			// execute return template which hope is the leaf template
			var updated_target common.Target
			err = yaml.Unmarshal(buf.Bytes(), &updated_target)
			if err != nil {
				log.Fatalf("Unmarshal2: %v", err)
			}
			// unmarshal the update target to create the execution path
			tokens := strings.Split(updated_target.Action, "/")
			repo_path, ok := lookUpRepoPath[tokens[0]]
			if !ok {
				log.Fatalln("Repo prefix does not exist!")
			}
			result = strings.TrimSpace(action.RunAction(updated_target,path.Join(tokens[1:]...),&action.ActionArgs{
				RepoPath: repo_path,
				TempDir: temp_dir,
				SpecMap: map[string]func(*action.ActionArgs) string{
					action.CLISpec: action.ExecCLI,
					action.WrapSpec: action.RunWrap,
				},
			}))
			lookUpReturn[string(yamlTarget)] = result
		}
		return result
	}
}

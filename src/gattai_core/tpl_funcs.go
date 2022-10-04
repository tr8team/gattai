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

func TplFetch(gattai_file GattaiFile, temp_dir string, lookUpRepoPath map[string]string, lookUpReturn map[string]string,specMap map[string]action.ActionFunc) func(common.Target) string {
	return func(target common.Target) string {
		// get target generated key
		yamlTarget, err := yaml.Marshal(target)
		if err != nil {
			log.Fatalf("TplFetch Marshal error: %v", err)
		}
		var buf bytes.Buffer
		// check if result for target already exist
		result, ok := lookUpReturn[string(yamlTarget)]
		if !ok {
			// if not, parse target to see if target have dependency
			tmpl, err := template.New("").Funcs(template.FuncMap{
				"fetch": TplFetch(gattai_file,temp_dir,lookUpRepoPath,lookUpReturn,specMap),
			}).Parse(string(yamlTarget))
			if err != nil {
				log.Fatalf("TplFetch template Parse error: %v", err)
			}
			buf.Reset()
			err = tmpl.Execute(&buf, gattai_file);
			if err != nil {
				log.Fatalf("TplFetch Execute error: %v", err)
			}
			// execute return template which hope is the leaf template
			var updated_target common.Target
			err = yaml.Unmarshal(buf.Bytes(), &updated_target)
			if err != nil {
				log.Fatalf("TplFetch Unmarshal error: %s error: %v", buf.String(), err)
			}
			// unmarshal the update target to create the execution path
			tokens := strings.Split(updated_target.Action, "/")
			repo_path, ok := lookUpRepoPath[tokens[0]]
			if !ok {
				log.Fatalln("TplFetch lookUpRepoPath error")
			}
			output, err := action.RunAction(updated_target,path.Join(tokens[1:]...),&action.ActionArgs{
				RepoPath: repo_path,
				TempDir: temp_dir,
				SpecMap: specMap,
			})
			if err != nil {
				log.Fatalf("TplFetch RunAction error: %v", err)
			}
			result = strings.TrimSpace(output)
			lookUpReturn[string(yamlTarget)] = result
		}
		return result
	}
}

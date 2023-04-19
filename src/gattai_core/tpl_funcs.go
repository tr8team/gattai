package core

import (
	//"os"
	//"io"
	//"fmt"
	"log"
	//"time"
	"path"
	"sync"
	"bytes"
	"strings"
	//"context"
	//"runtime"
	//"sync/atomic"
	"text/template"
	"gopkg.in/yaml.v2"
	"github.com/tr8team/gattai/src/gattai_core/action"
	"github.com/tr8team/gattai/src/gattai_core/common"
)

type LookUp struct {
	m sync.Map
}

func MakeLookUp() LookUp {
	return LookUp {}
}

func (lu *LookUp) Set(key string, val string) {
    lu.m.Store(key, val)
}

func (lu *LookUp) Get(key string) (string,bool) {
	if val, ok := lu.m.Load(key); ok {
        return val.(string), true
    }
    return "", false
}

func GoroutineFetch(yamlTarget string, gattai_file GattaiFile, temp_dir string, lookUpRepoPath map[string]string, lookUpReturn LookUp,cmdFunc action.CommandFunc, output chan string) {
	result, ok := lookUpReturn.Get(yamlTarget)
	if !ok {
		// if not, parse target to see if target have dependency
		tmpl, err := template.New("").Funcs(template.FuncMap{
			"fetch": TplFetch(gattai_file,temp_dir,lookUpRepoPath,lookUpReturn,cmdFunc),
		}).Parse(yamlTarget)
		if err != nil {
			log.Fatalf("TplFetch template Parse error: %v", err)
		}
		var buf bytes.Buffer
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
		tmpl_filepath := path.Join(repo_path,path.Join(tokens[1:]...)) + ".yaml"
		act_args := action.ActionArgs{
			RepoPath: repo_path,
			TempDir: temp_dir,
			SpecMap: map[string]action.ActionFunc{
				action.ActionVerKey(action.CLISpec, action.Version1): action.NewSpec[action.CommandLineInteraceSpec],
				action.ActionVerKey(action.DerivedSpec, action.Version1): action.NewSpec[action.DerivedInterfaceSpec],
			},
		}
		out_spec, err := action.RunAction(updated_target,tmpl_filepath,act_args)
		if err != nil {
			log.Fatalf("TplFetch RunAction error: %v", err)
		}
		out_result, err := cmdFunc(out_spec,act_args,updated_target.Action)
		if err != nil {
			log.Fatalf("TplFetch cmdFunc error: %v", err)
		}
		result = strings.TrimSpace(out_result)
		lookUpReturn.Set(string(yamlTarget), result)
	}
	output <- result
}

func TplFetch(gattai_file GattaiFile, temp_dir string, lookUpRepoPath map[string]string, lookUpReturn LookUp,cmdFunc action.CommandFunc) func(common.Target) string {
	return func(target common.Target) string {
		// get target generated key
		yamlTarget, err := yaml.Marshal(target)
		if err != nil {
			log.Fatalf("TplFetch Marshal error: %v", err)
		}
		// check if result for target already exist
		result := make(chan string)
		go GoroutineFetch(string(yamlTarget), gattai_file, temp_dir, lookUpRepoPath, lookUpReturn,cmdFunc, result)
		return <- result
	}
}

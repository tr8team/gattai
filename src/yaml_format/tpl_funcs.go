package yaml_format

import (
	"log"
	"path"
	"bytes"
	"strings"
	"text/template"
	"gopkg.in/yaml.v2"
	"github.com/tr8team/gattai/src/gattai_core/core_engine"
)

func YamlTarget(gattai_file GattaiFile, temp_dir string, lookUpRepoPath map[string]string,cmdFn CommandFunc) func(string) core_engine.FetchFunc{
	return func(yamlTargetBody string) core_engine.FetchFunc {
		return func(engine *core_engine.Engine) (string,error) {
			// if not, parse target to see if target have dependency
			tmpl, err := template.New("").Funcs(template.FuncMap{
				"fetch": TplFetch(gattai_file,temp_dir,lookUpRepoPath,engine,cmdFn),
			}).Parse(yamlTargetBody)
			if err != nil {
				log.Fatalf("YamlTarget template Parse error: %v", err)
			}
			var buf bytes.Buffer
			err = tmpl.Execute(&buf, gattai_file);
			if err != nil {
				log.Fatalf("YamlTarget Execute error: %v", err)
			}
			// execute return template which hope is the leaf template
			var updated_target Target
			err = yaml.Unmarshal(buf.Bytes(), &updated_target)
			if err != nil {
				log.Fatalf("YamlTarget Unmarshal error: %s error: %v", buf.String(), err)
			}
			// unmarshal the update target to create the execution path
			tokens := strings.Split(updated_target.Action, "/")
			repo_path, ok := lookUpRepoPath[tokens[0]]
			if !ok {
				log.Fatalln("YamlTarget lookUpRepoPath error")
			}
			tmpl_filepath := path.Join(repo_path,path.Join(tokens[1:]...)) + ".yaml"
			act_args := ActionArgs{
				RepoPath: repo_path,
				TempDir: temp_dir,
				SpecMap: map[string]ActionFunc{
					ActionVerKey(CLISpec, ActionVersion1): NewSpec[CommandLineInteraceSpec],
					ActionVerKey(DerivedSpec, ActionVersion1): NewSpec[DerivedInterfaceSpec],
				},
			}
			out_spec, err := RunAction(updated_target,tmpl_filepath,act_args)
			if err != nil {
				log.Fatalf("TplFetch RunAction error: %v", err)
			}
			return cmdFn(out_spec,act_args,updated_target.Action)
		}
	}
}

func TplFetch(gattai_file GattaiFile, temp_dir string, lookUpRepoPath map[string]string, engine *core_engine.Engine,cmdFn CommandFunc) func(Target) string {
	return func(target Target) string {
		// get target generated key
		yamlTarget, err := yaml.Marshal(target)
		if err != nil {
			log.Fatalf("TplFetch Marshal error: %v", err)
		}
		return engine.Fetch(string(yamlTarget),YamlTarget(gattai_file, temp_dir, lookUpRepoPath, cmdFn)(string(yamlTarget)))
	}
}

func TplFormat() func(string) string {
	return func(content string) string {
		new_content := strings.ReplaceAll(content, "\"","\\\"")
		return strings.ReplaceAll(new_content, "\n","\\n")
	}
}

func TplTempDir(temp_dir string) func(string) string {
	return func(filename string) string {
		return path.Join(temp_dir,filename)
	}
}

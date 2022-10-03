package action

import (
	//"fmt"
	"log"
	"path"
	"bytes"
	"text/template"
	"gopkg.in/yaml.v2"
	"github.com/tr8team/gattai/src/gattai_core/common"
)

const (
	Version1 string = "v1"
)

const (
	CLISpec string = "CommandLineInterface"
	WrapSpec       = "WrapperInterface"
)

const (
	StrBool string 	= "bool"
	StrInt        	= "int"
	StrFlt        	= "float"
	StrStr       	= "string"
	StrArr        	= "array"
	StrObj       	= "object"
)

type ActionFile struct {
	Version string `yaml:"version"`
	Type string `yaml:"type"`
	Params Params `yaml:"params"`
	Spec interface{} `yaml:"spec"`
}

type Param struct {
	Desc string `yaml:"desc"`
	Type string `yaml:"type"`
	Properties Params `yaml:"properties"`
}

type Params struct {
	Required map[string]*Param `yaml:"required"`
	Optional map[string]*Param `yaml:"optional"`
}

type ActionFunc func(common.Target,ActionFile,*ActionArgs)string

type ActionArgs struct {
	RepoPath string
	TempDir string
	SpecMap map[string]ActionFunc
}

func ActionVerKey(action string, ver string) string {
	return action + ver
}

func ValType(item interface{}) string {

	switch i_type := item.(type) {
	case bool:
		return StrBool
	case int:
		return StrInt
	case int8:
		return StrInt
	case int16:
		return StrInt
	case int32:
		return StrInt
	case int64:
		return StrInt
	case uint:
		return StrInt
	case uint8:
		return StrInt
	case uint16:
		return StrInt
	case uint32:
		return StrInt
	case uint64:
		return StrInt
	case float32:
		return StrFlt
	case float64:
		return StrFlt
	case string:
		return StrStr
	case []interface{}:
		return StrArr
	case map[interface{}]interface{}:
		return StrObj
	default:
		log.Fatalf("Unsupported type: %T!\n", i_type)
	}
	return ""
}

func  NewSpec[T any](actionFile ActionFile) *T {
	newSpec := new(T)
	yamlSpec, err := yaml.Marshal(actionFile.Spec)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s\n", yamlSpec)
	err = yaml.Unmarshal(yamlSpec, newSpec)
	if err != nil {
		log.Fatalf("Unmarshal4 newSpec: %v", err)
	}
	return newSpec
}

func (actionFile ActionFile) CheckParams(target common.Target) {
	check_params_rec(target, actionFile.Params)
}

func check_params_rec(target common.Target,params Params) {
	for key, val := range params.Required{
		if var_item, ok := target.Vars[key]; ok {
			var_type := ValType(var_item)
			if val.Type == var_type {
				if val.Type == StrObj {
					check_params_rec(target,val.Properties)
				}
			} else {
				log.Fatalf("Invalid type for %s: %v, Expecting %v",key,var_type,val.Type)
			}
		} else {
			log.Fatalf("Missing key %s, key is required!",key)
		}
	}
	for key, val := range params.Optional {
		if var_item, ok := target.Vars[key]; ok {
			var_type := ValType(var_item)
			if val.Type == var_type {
				if val.Type == StrObj {
					check_params_rec(target,val.Properties)
				}
			} else {
				log.Fatalf("Invalid type for %s: %v, Expecting %v",key,var_type,val.Type)
			}
		}
	}
}

func RunAction(updated_target common.Target, exec_filename string, action_args *ActionArgs) string {

	tmpl_filepath := path.Join(action_args.RepoPath,exec_filename) + ".yaml"
	tmpl_filename := path.Base(tmpl_filepath)
	tmpl, err := template.New(tmpl_filename).Funcs(template.FuncMap{
		"temp_dir": TplTempDir(action_args.TempDir),
		"format": TplFormat(),
	}).ParseFiles(tmpl_filepath)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, updated_target); err != nil {
		panic(err)
	}
	var actionFile ActionFile;

	//fmt.Printf("%s\n", buf.String())
	err = yaml.Unmarshal(buf.Bytes(), &actionFile)
	if err != nil {
		log.Fatalf("Unmarshal3: %v", err)
	}

	switch actionFile.Version {
	case Version1:
	default:
		log.Fatalf("This version is not supported: %v", actionFile.Version)
	}

	actionFile.CheckParams(updated_target)
	var result string

	if spec, ok := action_args.SpecMap[ActionVerKey(actionFile.Type,actionFile.Version)]; ok {
		result += spec(updated_target,actionFile,action_args)
	} else {
		log.Fatalf("Action file type with version is not supported: %s %s!", actionFile.Type,actionFile.Version)
	}

	return result
}

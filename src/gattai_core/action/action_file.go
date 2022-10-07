package action

import (
	"fmt"
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
	DerivedSpec       = "DerivedInterface"
)

const (
	StrBool string 	= "bool"
	StrInt        	= "int"
	StrFlt        	= "float"
	StrStr       	= "string"
	StrObj       	= "object"
	StrList        	= "list"
	StrDict			= "dict"
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
	ObjectOf Params `yaml:"object_of"`
	ListOf *Param `yaml:"list_of"`
	DictOf struct {
		Key_Type string `yaml:"key_type"`
		Value *Param `yaml:"value"`
	}`yaml:"dict_of"`
}

type Params struct {
	Required map[string]*Param `yaml:"required"`
	Optional map[string]*Param `yaml:"optional"`
}

type ActionFunc func(common.Target,ActionFile,*ActionArgs) (string,error)

type ActionArgs struct {
	RepoPath string
	TempDir string
	SpecMap map[string]ActionFunc
}

func ValPlainType(item interface{}) (string,error) {
	var result string
	switch i_type := item.(type) {
	case bool:
		result = StrBool
	case int:
		result = StrInt
	case int8:
		result = StrInt
	case int16:
		result = StrInt
	case int32:
		result = StrInt
	case int64:
		result = StrInt
	case uint:
		result = StrInt
	case uint8:
		result = StrInt
	case uint16:
		result = StrInt
	case uint32:
		result = StrInt
	case uint64:
		result = StrInt
	case float32:
		result = StrFlt
	case float64:
		result = StrFlt
	case string:
		result = StrStr
	default:
		return result, fmt.Errorf("ValPlainType invalid error: %T",i_type)
	}
	return result, nil
}

func NewSpecFromBuffer[T any](buffer []byte) (*T,error) {
	newSpec := new(T)
	err := yaml.Unmarshal(buffer, newSpec)
	if err != nil {
		return new(T), fmt.Errorf("NewSpecFromBuffer Unmarshal error: %s error: %v",string(buffer),err)
	}
	return newSpec, nil
}

func  NewSpec[T any](actionFile ActionFile) (*T,error) {
	yamlSpec, err := yaml.Marshal(actionFile.Spec)
	if err != nil {
		return new(T), fmt.Errorf("NewSpecFromBuffer Marshal error: %v",err)
	}
	return NewSpecFromBuffer[T](yamlSpec)
}

func ActionVerKey(action string, ver string) string {
	return action + ver
}

func (actionFile ActionFile) CheckVersion() error {
	switch actionFile.Version {
	case Version1:
	default:
		return fmt.Errorf("ActionFile:CheckVersion inalid version error: %s",actionFile.Version)
	}
	return nil
}

func (actionFile ActionFile) CheckParams(target common.Target) error {
	switch var_item_type := target.Vars.(type) {
	case map[interface{}]interface{}:
		result, err := check_multi_params(var_item_type, actionFile.Params)
		if err != nil {
			return fmt.Errorf("ActionFile:CheckParams error: %v",err)
		}
		if len(result) > 0 {
			return fmt.Errorf("ActionFile:CheckParams error: %s",result)
		}
	default:
		return fmt.Errorf("ActionFile:CheckParams error: got %T expecting object!",var_item_type)
	}
	return nil
}

func check_plain_type(var_item interface{}, val_type string) (string,error){
	var result string
	var_type, err := ValPlainType(var_item)
	if err != nil {
		return result, fmt.Errorf("check_single_param error: %v",err)
	}
	if var_type != val_type {
		result += fmt.Sprintf("check_single_param invalid type error: got %s expecting %s\n",var_type,val_type)
	}
	return result, nil
}

func check_single_param(var_item interface{},val *Param) (string,error){
	var result string
	switch var_item_type := var_item.(type) {
	case map[interface{}]interface{}:
		switch val.Type {
		case StrObj:
			output, err := check_multi_params(var_item_type,val.ObjectOf)
			if err != nil {
				return result, fmt.Errorf("check_single_param error: %v",err)
			}
			result += output
		case StrDict:
			for var_key, var_item := range var_item_type {
				output, err := check_plain_type(var_key,val.DictOf.Key_Type)
				if err != nil {
					return result, fmt.Errorf("check_single_param error: %v",err)
				}
				result += output
				output, err = check_single_param(var_item,val.DictOf.Value)
				if err != nil {
					return result, fmt.Errorf("check_single_param error: %v",err)
				}
				result += output
			}
		default:
			result += fmt.Sprintf("check_single_param invalid type error: map expecting %s\n",val.Type)
		}
	case []interface{}:
		if val.Type == StrList {
			for _, var_item := range var_item_type {
				output, err := check_single_param(var_item,val.ListOf)
				if err != nil {
					return result, fmt.Errorf("check_single_param error: %v",err)
				}
				result += output
			}
		}
	default:
		output, err := check_plain_type(var_item,val.Type)
		if err != nil {
			return result, fmt.Errorf("check_single_param error: %v",err)
		}
		result += output
	}

	return result,nil
}

func check_multi_params(target_var map[interface{}]interface{},params Params) (string,error) {
	var result string
	for key, val := range params.Required{
		if var_item, ok := target_var[key]; ok {
			output, err := check_single_param(var_item,val)
			if err != nil {
				return result, fmt.Errorf("check_multi_params error: %v",err)
			}
			result += output
		} else {
			result += fmt.Sprintf("check_multi_params:%s key is required error\n",key)
		}
	}
	for key, val := range params.Optional{
		if var_item, ok := target_var[key]; ok {
			output, err := check_single_param(var_item,val)
			if err != nil {
				return result, fmt.Errorf("check_multi_params error: %v",err)
			}
			result += output
		}
	}
	return result, nil
}

func RunAction(updated_target common.Target, tmpl_filepath string, action_args *ActionArgs) (string,error) {
	var result string

	tmpl_filename := path.Base(tmpl_filepath)
	tmpl, err := template.New(tmpl_filename).Funcs(template.FuncMap{
		"temp_dir": TplTempDir(action_args.TempDir),
		"format": TplFormat(),
	}).ParseFiles(tmpl_filepath)
	if err != nil {
		return result, fmt.Errorf("RunAction template error: %v",err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, updated_target)
	if err != nil {
		return result, fmt.Errorf("RunAction Execute error: %v",err)
	}
	var actionFile ActionFile;

	err = yaml.Unmarshal(buf.Bytes(), &actionFile)
	if err != nil {
		return result, fmt.Errorf("RunAction Unmarshal error: %s error: %v",buf.String(),err)
	}

	err = actionFile.CheckVersion()
	if err != nil {
		return result, fmt.Errorf("RunAction CheckVersion error: %v",err)
	}

	err = actionFile.CheckParams(updated_target)
	if err != nil {
		return result, fmt.Errorf("RunAction CheckParams error: %v",err)
	}

	if spec, ok := action_args.SpecMap[ActionVerKey(actionFile.Type,actionFile.Version)]; ok {
		output, err := spec(updated_target,actionFile,action_args)
		if err != nil {
			return result, fmt.Errorf("RunAction spec error: %v",err)
		}
		result += output
	} else {
		return result, fmt.Errorf("RunAction action type and version error: %s %s", actionFile.Type,actionFile.Version)
	}

	return result, nil
}

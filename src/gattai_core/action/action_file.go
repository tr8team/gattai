package action

import (
	"fmt"
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
	//Deprecated string `yaml:"deprecated"`
	Spec interface{} `yaml:"spec"`
}

type Param struct {
	Desc string `yaml:"desc"`
	Type string `yaml:"type"`
	ObjectOf Params `yaml:"object_of"`
	ListOf *Param `yaml:"list_of"`
	DictOf struct {
		Key *Param  `yaml:"key"`
		Value *Param `yaml:"value"`
	}`yaml:"dict_of"`
}

type Params struct {
	Required map[string]*Param `yaml:"required"`
	Optional map[string]*Param `yaml:"optional"`
}

type CommandFunc func(ActionSpec,ActionArgs,string) (string,error)

type ActionFunc func(ActionFile)(ActionSpec,error)

type ActionArgs struct {
	RepoPath string
	TempDir string
	SpecMap map[string]ActionFunc
}

type ParamField struct {
	Name string  `yaml:"Name"`
	Desc string `yaml:"Desc"`
	Attribute string `yaml:"Attribute"`
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

func NewSpecFromBuffer[T ActionSpec](buffer []byte) (ActionSpec,error) {
	var newSpec T
	err := yaml.Unmarshal(buffer, &newSpec)
	if err != nil {
		return newSpec, fmt.Errorf("NewSpecFromBuffer Unmarshal error: %s error: %v",string(buffer),err)
	}
	return newSpec, nil
}

func  NewSpec[T ActionSpec](actionFile ActionFile) (ActionSpec,error) {
	yamlSpec, err := yaml.Marshal(actionFile.Spec)
	if err != nil {
		var newSpec T
		return newSpec, fmt.Errorf("NewSpecFromBuffer Marshal error: %v",err)
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

func (actionFile ActionFile) GenerateParamFields() []ParamField {
	result := rec_paramfield_multi_params(actionFile.Params)
	log.Printf("GenerateParamFields:\n%v\n",result)
	return result
}

func rec_paramfield_multi_params(params Params) []ParamField {
	result := []ParamField{}
	for key, val := range params.Required{
		result = append(result, rec_paramfield_single_param(val, key, "<b>(required)</b>")...)
	}
	for key, val := range params.Optional{
		result = append(result, rec_paramfield_single_param(val, key, "<i>(optional)</i>")...)
	}
	return result
}

func rec_paramfield_single_param(val *Param,key string, attrib string) []ParamField {
	result := []ParamField{}
	switch val.Type {
	case StrObj:
		result = append(result, rec_paramfield_multi_params(val.ObjectOf)...)
	case StrList:
		result = append(result, ParamField{
			Name: key,
			Desc: val.Desc,
			Attribute: attrib,
		})
		if val.ListOf.Type == StrObj {
			result = append(result, rec_paramfield_multi_params(val.ListOf.ObjectOf)...)
		}
	case StrDict:
		result = append(result, ParamField{
			Name: key,
			Desc: val.Desc,
			Attribute: attrib,
		})
		if val.DictOf.Value.Type == StrObj {
			result = append(result, rec_paramfield_multi_params(val.DictOf.Value.ObjectOf)...)
		}
	default:
		result = append(result, ParamField{
			Name: key,
			Desc: val.Desc,
			Attribute: attrib,
		})
	}
	return result
}

func (actionFile ActionFile) GenerateTargetFromParamsInYaml(filename string) (string,error) {
	vars := make(map[interface{}]interface{})
	err := rec_target_from_multi_params(actionFile.Params,&vars)
	if err != nil {
		return "", fmt.Errorf("GenerateTargetFromParamsInYaml error: %v",err)
	}
	var extension = path.Ext(filename)
	var action_name = filename[0:len(filename)-len(extension)]
	result := map[string]interface{}{
		"action": path.Join("<repo_id>",action_name),
		"vars": &vars,
	}
	yamlFmt, err := yaml.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("rec_target_from_multi_params Marshal error: %v",err)
	}
	return string(yamlFmt), nil
}

func rec_target_from_multi_params(params Params,out *map[interface{}]interface{}) error {
	for key, val := range params.Required{
		var result interface{}
		err := rec_target_from_single_param(val,&result)
		if err != nil {
			return fmt.Errorf("rec_target_from_multi_params error: %v",err)
		}
		(*out)[key] = result
	}
	for key, val := range params.Optional{
		var result interface{}
		err := rec_target_from_single_param(val,&result)
		if err != nil {
			return fmt.Errorf("rec_target_from_multi_params error: %v",err)
		}
		(*out)[key] = result
	}
	return nil
}

func rec_target_from_single_param(val *Param,out *interface{}) error {
	switch val.Type {
	case StrObj:
		result := make(map[interface{}]interface{})
		err := rec_target_from_multi_params(val.ObjectOf,&result)
		if err != nil {
			return fmt.Errorf("rec_target_from_single_param error: %v",err)
		}
		*out = &result
	case StrList:
		var result interface{}
		err := rec_target_from_single_param(val.ListOf,&result)
		if err != nil {
			return fmt.Errorf("rec_target_from_single_param error: %v",err)
		}
		*out = &[]interface{}{result}
	case StrDict:
		var key interface{}
		err := rec_target_from_single_param(val.DictOf.Key,&key)
		if err != nil {
			return fmt.Errorf("rec_target_from_single_param error: %v",err)
		}
		var value interface{}
		err = rec_target_from_single_param(val.DictOf.Value,&value)
		if err != nil {
			return fmt.Errorf("rec_target_from_single_param error: %v",err)
		}
		*out = &map[interface{}]interface{}{key : value}
	default:
		result := "\"" + val.Type + "\""
		*out = &result
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
				output, err := check_plain_type(var_key,val.DictOf.Key.Type)
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

func RunAction(updated_target common.Target, tmpl_filepath string, action_args ActionArgs) (ActionSpec,error) {
	tmpl_filename := path.Base(tmpl_filepath)
	tmpl, err := template.New(tmpl_filename).Funcs(template.FuncMap{
		"temp_dir": TplTempDir(action_args.TempDir),
		"format": TplFormat(),
	}).ParseFiles(tmpl_filepath)
	if err != nil {
		return nil, fmt.Errorf("RunAction template error: %v",err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, updated_target)
	if err != nil {
		return nil, fmt.Errorf("RunAction Execute error: %v",err)
	}
	var actionFile ActionFile;

	err = yaml.Unmarshal(buf.Bytes(), &actionFile)
	if err != nil {
		return nil, fmt.Errorf("RunAction Unmarshal error: %s error: %v",buf.String(),err)
	}

	err = actionFile.CheckVersion()
	if err != nil {
		return nil, fmt.Errorf("RunAction CheckVersion error: %v",err)
	}

	err = actionFile.CheckParams(updated_target)
	if err != nil {
		return nil, fmt.Errorf("RunAction CheckParams error: %v",err)
	}

	if spec, ok := action_args.SpecMap[ActionVerKey(actionFile.Type,actionFile.Version)]; ok {
		return spec(actionFile)
	}

	return nil, fmt.Errorf("RunAction action type and version error: %s %s", actionFile.Type,actionFile.Version)
}

package action

import (
	"log"
	"github.com/tr8team/gattai/src/gattai_core/common"
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

func CheckParams(target common.Target,param_map map[string]*Param) {
	for key, val := range param_map {
		if var_item, ok := target.Vars[key]; ok {
			var_type := ValType(var_item)
			if val.Type == var_type {
				if val.Type == StrObj {
					CheckParams(target,val.Properties.Required)
				}
			} else {
				log.Fatalf("Invalid type for %s: %v, Expecting %v",key,var_type,val.Type)
			}
		} else {
			log.Fatalf("Missing key %s, key is required!",key)
		}
	}
}

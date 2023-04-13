package action

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

type ParamField struct {
	Name string  `yaml:"Name"`
	Desc string `yaml:"Desc"`
	Attribute string `yaml:"Attribute"`
}

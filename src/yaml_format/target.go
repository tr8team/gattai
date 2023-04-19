package yaml_format

type Target struct {
	Action string `yaml:"action"`
	Vars interface{} `yaml:"vars"`
}

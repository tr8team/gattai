package common

type Target struct {
	Action string `yaml:"action"`
	Vars interface{} `yaml:"vars"`
}

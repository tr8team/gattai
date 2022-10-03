package common

type Target struct {
	Action string `yaml:"action"`
	Vars map[string]interface{} `yaml:"vars"`
}

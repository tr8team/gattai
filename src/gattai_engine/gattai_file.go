package core

import (
	"github.com/tr8team/gattai/src/gattai_engine/common"
)

const (
	Version1 string = "v1"
)

const (
	AllNamespaces string = "all"
	AllTargets string = "all"
)

const (
	GattaiFileDefault string = "Gattaifile.yaml"
	GattaiTmpFolder string = "gattaitmp"
)

type GattaiFile struct {
    Version string `yaml:"version"`
    TempFolder string `yaml:"temp_folder"`
	EnforceTargets map[string][]string `yaml:"enforce_targets"`
	Repos map[string]common.Repo `yaml:"repos"`
	Targets map[string]map[string]common.Target `yaml:"targets"`
}

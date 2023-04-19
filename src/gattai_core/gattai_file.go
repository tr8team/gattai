package core

import (
	"os"
	"fmt"
	"errors"
	"gopkg.in/yaml.v2"
	"github.com/tr8team/gattai/src/yaml_format"
	"github.com/tr8team/gattai/src/gattai_core/common"
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

func NewGattaiFileFromBuffer(buffer []byte) (*GattaiFile,error) {
	gattaiFile := new(GattaiFile)
	err := yaml.Unmarshal(buffer, gattaiFile)
	if err != nil {
		return &GattaiFile{}, fmt.Errorf("NewGattaiFileFromBuffer Unmarshal error: %s error: %v",string(buffer),err)
	}
	return gattaiFile, nil
}

func NewGattaiFile(gattaifile_path string) (*GattaiFile,error) {
	yamlFile, err := os.ReadFile(gattaifile_path)
	if err != nil {
		return &GattaiFile{}, fmt.Errorf("NewGattaiFile osReadFile error: %s error: %v",gattaifile_path,err)
	}
	return NewGattaiFileFromBuffer(yamlFile)
}

func (gattaiFile GattaiFile) CheckVersion() error {
	switch gattaiFile.Version {
	case Version1:
	default:
		return fmt.Errorf("GattaiFile:CheckVersion invalid version error: %s",gattaiFile.Version)
	}
	return nil
}

func (gattaiFile GattaiFile) CheckEnforceTargets() error {
	result := ""
	for namespace_id, target_id_list := range gattaiFile.EnforceTargets {
		if targets, ok := gattaiFile.Targets[namespace_id]; ok {
			for _, target_id := range target_id_list {
				if _, ok := targets[target_id]; !ok {
					result += fmt.Sprintf("GattaiFile:CheckEnforceTargets:%s:%s failed to enforce error!\n",namespace_id, target_id)
				}
			}
		} else {
			for _, target_id := range target_id_list {
				result += fmt.Sprintf("GattaiFile:CheckEnforceTargets:%s:%s failed to enforce error!\n",namespace_id, target_id)
			}
		}
	}
	if len(result) > 0 {
		return errors.New(result)
	}
	return nil
}

func (gattaiFile GattaiFile) CreateTempDir(folder_prefix string) (string, error) {
	tempDir, err := os.MkdirTemp(gattaiFile.TempFolder, folder_prefix)
	if err != nil {
		return tempDir, fmt.Errorf("GattaiFile:CreateTempDir osMkdirTemp error: %s %s error: %v",gattaiFile.TempFolder,folder_prefix,err)
	}
	return tempDir, nil
}

func (gattaiFile GattaiFile) BuildRepoMap(tempDir string) (map[string]string, error) {
	result := make(map[string]string)

	for key, val := range gattaiFile.Repos {
		output, err  := common.GetRepoPath(tempDir,key,val,"")
		if err != nil {
			return result, fmt.Errorf("GattaiFile:BuildRepoMap error %v",err)
		}
		result[key] = output
	}

	return result, nil
}

func (gattaiFile GattaiFile) LookupTargets(namespace_id string, target_id string, tempDir string,cmdFunc yaml_format.CommandFunc) (string,error) {
	var result string

	lookUpRepoPath,err := gattaiFile.BuildRepoMap(tempDir)
	if err != nil {
		return result, fmt.Errorf("GattaiFile:LookupTargets error: %v",err)
	}
	lookUpReturn := MakeLookUp()

	switch namespace_id {
	case AllNamespaces:
		switch  target_id {
		case AllTargets:
			// all namespaces and all targets
			for _, targets := range gattaiFile.Targets {
				for _, target := range targets {
					result += TplFetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn,cmdFunc)(target)
				}
			}
		default:
			// all namespaces and a single target
			for _, targets := range gattaiFile.Targets {
				if target, ok := targets[target_id]; ok {
					result += TplFetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn,cmdFunc)(target)
				}
			}
		}
	default:
		if targets , ok := gattaiFile.Targets[namespace_id]; ok {
			switch  target_id {
			case AllTargets:
				// a single namespace and all targets
				for _, target := range targets {
					result += TplFetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn,cmdFunc)(target)
				}
			default:
				// a single namespace and a single target
				if target, ok := targets[target_id]; ok {
					result += TplFetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn,cmdFunc)(target)
				}
			}
		}
	}
	return result, nil
}

package core

import (
	"os"
	"fmt"
	"errors"
	"net/url"
	"gopkg.in/yaml.v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tr8team/gattai/src/gattai_core/action"
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
	GattaiFileDefault string = "GattaiFile.yaml"
	GattaiTmpFolder string = "gattaitmp"
)

type GattaiFile struct {
    Version string `yaml:"version"`
    TempFolder string `yaml:"temp_folder"`
	EnforceTargets map[string][]string `yaml:"enforce_targets"`
	Repos map[string]struct {
		Repo string `yaml:"repo"`
		Src map[string]string `yaml:"src"`
	} `yaml:"repos"`
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

func (gattaiFile GattaiFile) CheckEnforceTargets() error {
	result := ""
	for namespace_id, target_id_list := range gattaiFile.EnforceTargets {
		if targets, ok := gattaiFile.Targets[namespace_id]; ok {
			for _, target_id := range target_id_list {
				if _, ok := targets[target_id]; !ok {
					result += fmt.Sprintf("CheckEnforceTargets:%s:%s failed to enforce error!\n",namespace_id, target_id)
				}
			}
		} else {
			for _, target_id := range target_id_list {
				result += fmt.Sprintf("CheckEnforceTargets:%s:%s failed to enforce error!\n",namespace_id, target_id)
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
		return tempDir, fmt.Errorf("CreateTempDir osMkdirTemp error: %s %s error: %v",gattaiFile.TempFolder,folder_prefix,err)
	}
	return tempDir, nil
}

func (gattaiFile GattaiFile) BuildRepoMap() (map[string]string, error) {
	result := make(map[string]string)

	for key, val := range gattaiFile.Repos {
		src := val.Src
		switch val.Repo {
		case "local":
			dir, ok := src["dir"]
			if ok == false {
				return result, errors.New("BuildRepoMap:local error: dir is missing!")
			}
			fileInfo, err := os.Stat(dir)
			if err != nil {
				return result, fmt.Errorf("BuildRepoMap:local osStat error: %s error: %v",dir,err)
			}
			if fileInfo.IsDir() == false {
				return result, fmt.Errorf("BuildRepoMap:local IsDir error: %s is not a directory!",dir)
			}
			result[key] = dir
		case "git":
			web_url, ok := src["url"]
			if ok == false {
				return result, errors.New("BuildRepoMap:git error: url is missing!")
			}
			_, err := url.ParseRequestURI(web_url)
			if err != nil {
				return result, fmt.Errorf("BuildRepoMap:git ParseRequestURI error: %s error: %v",web_url,err)
			}
			repoDir, err := os.MkdirTemp("",key)
			if err != nil {
				return result, fmt.Errorf("BuildRepoMap:git osMkdirTemp error: %s error: %v",key,err)
			}
			var ref_name plumbing.ReferenceName
			if branch, ok := src["branch"]; ok {
				ref_name = plumbing.NewBranchReferenceName(branch)
			}
			if tag, ok := src["tag"]; ok {
				ref_name = plumbing.NewTagReferenceName(tag)
			}
			defer os.RemoveAll(repoDir) // clean up
			_, err = git.PlainClone(repoDir, false, &git.CloneOptions{
				URL:               web_url,
				Progress: 		   os.Stdout,
				ReferenceName:	   ref_name,
			})
			if err != nil {
				return result, fmt.Errorf("BuildRepoMap:git PlainClone error: %s error: %v",repoDir,err)
			}
			result[key] = repoDir
		default:
			return result, fmt.Errorf("BuildRepoMap Repo type error: %s is not supported!",val.Repo)
		}
	}

	return result, nil
}

func (gattaiFile GattaiFile) LookupTargets(namespace_id string, target_id string, tempDir string,specMap map[string]action.ActionFunc) (string,error) {
	var result string

	lookUpRepoPath,err := gattaiFile.BuildRepoMap()
	if err != nil {
		return result, fmt.Errorf("LookupTargets error: %v",err)
	}
	lookUpReturn := make(map[string]string)

	switch namespace_id {
	case AllNamespaces:
		switch  target_id {
		case AllTargets:
			// all namespaces and all targets
			for _, targets := range gattaiFile.Targets {
				for _, target := range targets {
					result += TplFetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn,specMap)(target)
				}
			}
		default:
			// all namespaces and a single target
			for _, targets := range gattaiFile.Targets {
				if target, ok := targets[target_id]; ok {
					result += TplFetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn,specMap)(target)
				}
			}
		}
	default:
		if targets , ok := gattaiFile.Targets[namespace_id]; ok {
			switch  target_id {
			case AllTargets:
				// a single namespace and all targets
				for _, target := range targets {
					result += TplFetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn,specMap)(target)
				}
			default:
				// a single namespace and a single target
				if target, ok := targets[target_id]; ok {
					result += TplFetch(gattaiFile,tempDir,lookUpRepoPath,lookUpReturn,specMap)(target)
				}
			}
		}
	}
	return result, nil
}

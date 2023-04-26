package yaml_format

import (
	"os"
	"fmt"
	"path"
	"errors"
	"net/url"
	"gopkg.in/yaml.v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tr8team/gattai/src/gattai_core/core_engine"
)

const (
	GattaiVersion1 string = "v1"
)

const (
	AllNamespaces string = "all"
	AllTargets string = "all"
)

const (
	GattaiFileDefault string = "Gattaifile.yaml"
	GattaiTmpFolder string = "gattaitmp"
)

const (
	LocalRepo string = "local"
	GitRepo       	 = "git"
)

type GattaiFile struct {
    Version string `yaml:"version"`
    TempFolder string `yaml:"temp_folder"`
	EnforceTargets map[string][]string `yaml:"enforce_targets"`
	Repos map[string]Repo `yaml:"repos"`
	Targets map[string]map[string]Target `yaml:"targets"`
}

type Repo struct {
	Src string `yaml:"src"`
	Config struct {
		Url string `yaml:"url"`
		Branch string `yaml:"branch"`
		Tag string `yaml:"tag"`
		Dir string `yaml:"dir"`
	} `yaml:"config"`
}

type Target struct {
	Action string `yaml:"action"`
	Vars interface{} `yaml:"vars"`
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
	case GattaiVersion1:
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


func GetRepoPath(tempDir string,key string,repo Repo,parent_repopath string) (string,error) {
	var result string
	switch repo.Src {
	case LocalRepo:
		repoDir := path.Join(parent_repopath,repo.Config.Dir)
		fileInfo, err := os.Stat(repoDir)
		if err != nil {
			return result, fmt.Errorf("GetRepoPath:local osStat error: %s error: %v",repoDir,err)
		}
		if fileInfo.IsDir() == false {
			return result, fmt.Errorf("GetRepoPath:local IsDir error: %s is not a directory!",repoDir)
		}
		result = repoDir
	case GitRepo:
		web_url := repo.Config.Url
		_, err := url.ParseRequestURI(web_url)
		if err != nil {
			return result, fmt.Errorf("GetRepoPath:git ParseRequestURI error: %s error: %v",web_url,err)
		}
		repoDir, err := os.MkdirTemp(tempDir,key)
		if err != nil {
			return result, fmt.Errorf("GetRepoPath:git osMkdirTemp error: %s error: %v",key,err)
		}
		var ref_name plumbing.ReferenceName
		branch := repo.Config.Branch
		if len(branch) > 0 {
			ref_name = plumbing.NewBranchReferenceName(branch)
		}
		tag := repo.Config.Tag;
		if len(tag) > 0 {
			ref_name = plumbing.NewTagReferenceName(tag)
		}
		_, err = git.PlainClone(repoDir, false, &git.CloneOptions{
			URL:               web_url,
			Progress: 		   os.Stdout,
			ReferenceName:	   ref_name,
		})
		if err != nil {
			return result, fmt.Errorf("GetRepoPath:git PlainClone error: %s error: %v",repoDir,err)
		}
		result = path.Join(repoDir,repo.Config.Dir)
	default:
		return result, fmt.Errorf("GetRepoPath Src type error: %s is not supported!",repo.Src)
	}
	return result, nil
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
		output, err  := GetRepoPath(tempDir,key,val,"")
		if err != nil {
			return result, fmt.Errorf("GattaiFile:BuildRepoMap error %v",err)
		}
		result[key] = output
	}

	return result, nil
}

func (gattaiFile GattaiFile) LookupTargets(namespace_id string, target_id string, tempDir string,cmdFn core_engine.CommandFunc) (string,error) {
	var result string

	lookUpRepoPath,err := gattaiFile.BuildRepoMap(tempDir)
	if err != nil {
		return result, fmt.Errorf("GattaiFile:LookupTargets error: %v",err)
	}
	engine := core_engine.MakeEngine(cmdFn)

	switch namespace_id {
	case AllNamespaces:
		switch  target_id {
		case AllTargets:
			// all namespaces and all targets
			for _, targets := range gattaiFile.Targets {
				for _, target := range targets {
					result += TplFetch(gattaiFile,tempDir,lookUpRepoPath,engine)(target)
				}
			}
		default:
			// all namespaces and a single target
			for _, targets := range gattaiFile.Targets {
				if target, ok := targets[target_id]; ok {
					result += TplFetch(gattaiFile,tempDir,lookUpRepoPath,engine)(target)
				}
			}
		}
	default:
		if targets , ok := gattaiFile.Targets[namespace_id]; ok {
			switch  target_id {
			case AllTargets:
				// a single namespace and all targets
				for _, target := range targets {
					result += TplFetch(gattaiFile,tempDir,lookUpRepoPath,engine)(target)
				}
			default:
				// a single namespace and a single target
				if target, ok := targets[target_id]; ok {
					result += TplFetch(gattaiFile,tempDir,lookUpRepoPath,engine)(target)
				}
			}
		}
	}
	return result, nil
}

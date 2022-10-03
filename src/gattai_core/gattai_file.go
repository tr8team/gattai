package core

import (
	"os"
	//"fmt"
	"log"
	"net/url"
	"io/ioutil"
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
	GattaiTmpFolder string = "gattaitmp"
)

const (
	RepoLocal string = "local"
	RepoGit          = "git"
)

const (
	LocalDir string = "dir"
)

const (
	GitUrl string = "url"
	GitTag        = "tag"
	GitBranch     = "branch"
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

func NewGattaiFile(gattaifile_path string) *GattaiFile {
	gattaiFile := new(GattaiFile)

	yamlFile, err := ioutil.ReadFile(gattaifile_path)
	if err != nil {
		log.Fatalf("Error reading Gattai File: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, gattaiFile)
	if err != nil {
		log.Fatalf("Error parsing Gattai File: %v", err)
	}

	return gattaiFile
}

func (gattaiFile GattaiFile) CheckEnforceTargets() {
	for namespace_id, target_id_list := range gattaiFile.EnforceTargets {
		if targets, ok := gattaiFile.Targets[namespace_id]; ok {
			for _, target_id := range target_id_list {
				if _, ok := targets[target_id]; !ok {
					log.Fatalf("Target from <%v> is required by enforced-target: %v", namespace_id, target_id)
				}
			}
		} else {
			log.Fatalf("Namespace is required by enforced-target: %v", namespace_id)
		}
	}
}

func (gattaiFile GattaiFile) CreateTempDir() string {
	tempDir, err := os.MkdirTemp(gattaiFile.TempFolder, GattaiTmpFolder)
	if err != nil {
		log.Fatalf("Error creating temporary folder: %v", err)
	}
	return tempDir
}

func (gattaiFile GattaiFile) BuildRepoMap() map[string]string {
	result := make(map[string]string)

	for key, val := range gattaiFile.Repos {
		src := val.Src
		switch val.Repo {
		case RepoLocal:
			dir, ok := src[LocalDir]
			if ok == false {
				log.Fatalln("Please provide a dir: path")
			}
			fileInfo, err := os.Stat(dir)
			if err != nil || fileInfo.IsDir() == false {
				log.Fatalln("Please provide a directory for local repo!")
			}
			result[key] = dir
		case RepoGit:
			web_url, ok := src[GitUrl]
			if ok == false {
			 log.Fatalln("Please provide a url: key")
			}
			_, err := url.ParseRequestURI(web_url)
			if err != nil {
				log.Fatalf("GIT repo parse request url error: %v", err)
			}
			repoDir, err := os.MkdirTemp("",key)
			if err != nil {
				log.Fatalf("Error creating repository folder: %v", err)
			}
			var ref_name plumbing.ReferenceName
			if branch, ok := src[GitBranch]; ok {
				ref_name = plumbing.NewBranchReferenceName(branch)
			}
			if tag, ok := src[GitTag]; ok {
				ref_name = plumbing.NewTagReferenceName(tag)
			}
			defer os.RemoveAll(repoDir) // clean up
			_, err = git.PlainClone(repoDir, false, &git.CloneOptions{
				URL:               web_url,
				Progress: 		   os.Stdout,
				ReferenceName:	   ref_name,
			})
			if err != nil {
				log.Fatalf("Error cloning git repository: %v", err)
			}
			result[key] = repoDir
		default:
			log.Fatalln("Repo type is not supported!")
		}
	}

	return result
}

func (gattaiFile GattaiFile) LookupTargets(namespace_id string, target_id string, tempDir string,specMap map[string]action.ActionFunc) string {
	var result string

	lookUpRepoPath := gattaiFile.BuildRepoMap()
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
	return result
}

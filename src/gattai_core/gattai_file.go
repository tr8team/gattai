package core

import (
	"os"
	//"io"
	//"fmt"
	"log"
	//"path"
	//"time"
	//"bytes"
	//"strconv"
	//"strings"
	//"context"
	//"runtime"
	//"os/exec"
	"net/url"
	//"io/ioutil"
	//"text/template"
	//"gopkg.in/yaml.v2"
	//"mvdan.cc/sh/v3/expand"
	//"mvdan.cc/sh/v3/interp"
	//"mvdan.cc/sh/v3/syntax"
	//"github.com/spf13/cobra"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tr8team/gattai/src/gattai_core/common"
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

func BuildRepoMap(gattaiFile GattaiFile) map[string]string {
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

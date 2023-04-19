package yaml_format

import (
	"os"
	"fmt"
	"path"
	//"errors"
	"net/url"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	LocalRepo string = "local"
	GitRepo       	 = "git"
)

type Repo struct {
	Src string `yaml:"src"`
	Config struct {
		Url string `yaml:"url"`
		Branch string `yaml:"branch"`
		Tag string `yaml:"tag"`
		Dir string `yaml:"dir"`
	} `yaml:"config"`
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

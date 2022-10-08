package common

import (
	"os"
	"fmt"
	"path"
	"errors"
	"net/url"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Repo struct {
	Src string `yaml:"src"`
	Config map[string]string `yaml:"config"`
}

func GetRepoPath(tempDir string,key string,repo Repo,parent_repopath string) (string,error) {
	var result string
	switch repo.Src {
	case "local":
		dir, ok := repo.Config["dir"]
		if ok == false {
			return result, errors.New("GetRepoPath:local error: dir is missing!")
		}
		repoDir := path.Join(parent_repopath,dir)
		fileInfo, err := os.Stat(dir)
		if err != nil {
			return result, fmt.Errorf("GetRepoPath:local osStat error: %s error: %v",dir,err)
		}
		if fileInfo.IsDir() == false {
			return result, fmt.Errorf("GetRepoPath:local IsDir error: %s is not a directory!",dir)
		}
		result = repoDir
	case "git":
		web_url, ok := repo.Config["url"]
		if ok == false {
			return result, errors.New("GetRepoPath:git error: url is missing!")
		}
		_, err := url.ParseRequestURI(web_url)
		if err != nil {
			return result, fmt.Errorf("GetRepoPath:git ParseRequestURI error: %s error: %v",web_url,err)
		}
		repoDir, err := os.MkdirTemp(tempDir,key)
		if err != nil {
			return result, fmt.Errorf("GetRepoPath:git osMkdirTemp error: %s error: %v",key,err)
		}
		var ref_name plumbing.ReferenceName
		if branch, ok := repo.Config["branch"]; ok {
			ref_name = plumbing.NewBranchReferenceName(branch)
		}
		if tag, ok := repo.Config["tag"]; ok {
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
		result = repoDir
	default:
		return result, fmt.Errorf("GetRepoPath Src type error: %s is not supported!",repo.Src)
	}
	return result, nil
}

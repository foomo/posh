package git

import (
	"github.com/go-git/go-git/v5"
	giturl "github.com/kubescape/go-git-url"
)

func OriginURL() (string, error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	if value, err := r.Remote("origin"); err != nil {
		return "", err
	} else if len(value.Config().URLs) == 0 {
		return "", nil
	} else if value, err := giturl.NewGitURL(value.Config().URLs[0]); err != nil {
		return "", err
	} else {
		return value.GetHostName() + "/" + value.GetRepoName(), nil
	}
}

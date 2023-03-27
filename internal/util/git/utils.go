package git

import (
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	giturls "github.com/whilp/git-urls"
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
	} else if value, err := giturls.Parse(value.Config().URLs[0]); err != nil {
		return "", err
	} else {
		return value.Hostname() + strings.TrimSuffix(value.Path, path.Ext(value.Path)), nil
	}
}

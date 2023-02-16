package files

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"

	"github.com/charlievieth/fastwalk"
)

type (
	FindOptions struct {
		ignore []*regexp.Regexp
		follow bool
	}
	FindOption func(*FindOptions)
)

func FindWithFollow(v bool) FindOption {
	return func(o *FindOptions) {
		o.follow = v
	}
}

func FindWithIgnore(v ...string) FindOption {
	return func(o *FindOptions) {
		for _, s := range v {
			o.ignore = append(o.ignore, regexp.MustCompile(s))
		}
	}
}

func Find(ctx context.Context, root, pattern string, opts ...FindOption) ([]string, error) {
	o := FindOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}
	var ret []string
	if err := fastwalk.Walk(&fastwalk.Config{
		Follow: o.follow,
	}, root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if err := ctx.Err(); err != nil {
			return err
		}

		if d.Name() != "." {
			for _, ignore := range o.ignore {
				if ignore.MatchString(d.Name()) {
					fmt.Println(d.Name())
					if d.IsDir() {
						return fastwalk.SkipDir
					} else {
						return nil
					}
				}
			}
		}

		if ok, err := filepath.Match(pattern, d.Name()); err != nil {
			return err
		} else if ok {
			ret = append(ret, p)
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

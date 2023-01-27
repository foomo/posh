package files

import (
	"context"
	"io/fs"
	"path/filepath"

	"github.com/charlievieth/fastwalk"
)

type (
	FindOptions struct {
		ignore []string
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
		o.ignore = append(o.ignore, v...)
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

		for _, pattern := range o.ignore {
			if ok, err := filepath.Match(pattern, d.Name()); err != nil {
				return err
			} else if ok {
				if d.IsDir() {
					return fastwalk.SkipDir
				} else {
					return nil
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

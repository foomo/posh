package files

import (
	"context"
	"io/fs"

	"github.com/charlievieth/fastwalk"
)

func Find(ctx context.Context, filename string) ([]string, error) {
	var ret []string
	if err := fastwalk.Walk(&fastwalk.Config{
		Follow: false,
	}, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if err := ctx.Err(); err != nil {
			return err
		}
		if d.Name() == filename {
			ret = append(ret, p)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

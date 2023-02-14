package files

import (
	"os"

	"github.com/pkg/errors"
)

func MkdirAll(paths ...string) error {
	for _, path := range paths {
		if path == "" {
			return errors.New("invalid empty path")
		}
		if stat, err := os.Stat(path); err != nil && os.IsNotExist(err) {
			if err := os.MkdirAll(path, os.ModeDir|0700); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else if !stat.IsDir() {
			return errors.Errorf("%s not a directory", path)
		}
	}
	return nil
}

package files

import (
	"os"

	"github.com/pkg/errors"
)

func Exists(paths ...string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			return errors.Wrap(err, path)
		}
	}
	return nil
}

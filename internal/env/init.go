package env

import (
	"os"

	"github.com/foomo/posh/pkg/env"
	"github.com/pkg/errors"
)

func Init() error {
	// setup env
	if value := os.Getenv(env.ProjectRoot); value != "" {
		// continue
	} else if wd, err := os.Getwd(); err != nil {
		return errors.Wrap(err, "failed to retrieve project root")
	} else if err := os.Setenv(env.ProjectRoot, wd); err != nil {
		return errors.Wrap(err, "failed to set project root env")
	}
	return nil
}

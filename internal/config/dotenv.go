package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func Dotenv() error {
	err := godotenv.Load(".poshrc")
	if errors.Is(err, os.ErrNotExist) {
		// continue
	} else if err != nil {
		return err
	}

	return nil
}

package config

import (
	"github.com/joho/godotenv"
)

func Dotenv() error {
	return godotenv.Load(".poshrc")
}

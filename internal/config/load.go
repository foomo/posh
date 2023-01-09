package config

import (
	"fmt"
	"os"

	"github.com/foomo/posh/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func Load(l log.Logger) error {
	var errNotFound viper.ConfigFileNotFoundError

	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".posh")
	if err := viper.ReadInConfig(); err != nil {
		return err
	} else {
		l.Debug("using config file:", viper.ConfigFileUsed())
	}

	override := viper.New()
	override.AddConfigPath(".")
	override.SetConfigType("yaml")
	override.SetConfigName(".posh.override")
	if err := override.ReadInConfig(); errors.As(err, &errNotFound) {
		// continue
	} else if err != nil {
		return err
	} else if err := viper.MergeConfigMap(override.AllSettings()); err != nil {
		return err
	} else {
		l.Debug("using override config file:", override.ConfigFileUsed())
	}

	// validate version
	if v := viper.GetString("version"); v != Version {
		return fmt.Errorf("invalid config version: %s (%s)", v, Version)
	}

	// set configured env
	for key, value := range viper.GetStringMapString("env") {
		if err := os.Setenv(key, os.ExpandEnv(value)); err != nil {
			return err
		}
	}

	return nil
}

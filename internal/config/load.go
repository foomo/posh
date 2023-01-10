package config

import (
	"fmt"
	"os"

	"github.com/foomo/posh/pkg/config"
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
	var env config.Env
	if err := viper.UnmarshalKey("env", &env); err != nil {
		l.Warn("failed to load env:", err.Error())
	} else {
		for _, value := range env {
			if err := os.Setenv(value.Name, os.ExpandEnv(value.Value)); err != nil {
				return err
			}
		}
	}

	return nil
}

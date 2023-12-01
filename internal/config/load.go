package config

import (
	"fmt"
	"os"

	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/log"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func Load(l log.Logger) error {
	var errNotFound viper.ConfigFileNotFoundError
	var settings map[string]interface{}

	if value := os.Getenv("POSH_ROOT_CONFIG_PATH"); value != "" {
		c := viper.New()
		c.AddConfigPath(value)
		c.SetConfigType("yaml")
		c.SetConfigName(".posh")
		if err := c.ReadInConfig(); errors.As(err, &errNotFound) {
			// continue
		} else if err != nil {
			return err
		} else if err := mergo.Merge(&settings, c.AllSettings(), mergo.WithAppendSlice); err != nil {
			return err
		} else {
			l.Debug("using root config file:", c.ConfigFileUsed())
		}
	}

	{ // load config
		c := viper.New()
		c.AddConfigPath(".")
		c.SetConfigType("yaml")
		c.SetConfigName(".posh")
		if err := c.ReadInConfig(); errors.As(err, &errNotFound) {
			// continue
		} else if err != nil {
			return err
		} else if err := mergo.Merge(&settings, c.AllSettings(), mergo.WithAppendSlice); err != nil {
			return err
		} else {
			l.Debug("using config file:", c.ConfigFileUsed())
		}
	}

	{ // load override
		c := viper.New()
		c.AddConfigPath(".")
		c.SetConfigType("yaml")
		c.SetConfigName(".posh.override")
		if err := c.ReadInConfig(); errors.As(err, &errNotFound) {
			// continue
		} else if err != nil {
			return err
		} else if err := mergo.Merge(&settings, c.AllSettings(), mergo.WithAppendSlice); err != nil {
			return err
		} else {
			l.Debug("using override config file:", c.ConfigFileUsed())
		}
	}

	if err := viper.MergeConfigMap(settings); err != nil {
		return errors.Wrap(err, "failed to merge config map")
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

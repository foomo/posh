package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/foomo/posh/pkg/log"
	"github.com/spf13/viper"
)

func Load(l log.Logger, configFile string) error {
	// setup viper
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".posh")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	} else {
		l.Debug("using config file:", viper.ConfigFileUsed())
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

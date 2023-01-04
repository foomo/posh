package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	EnvProjectRoot = "PROJECT_ROOT"
)

var (
	l           log.Logger
	m           *plugin.Manager
	flagLevel   string
	flagConfig  string
	flagNoColor bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "posh",
	Short: "Project Oriented Shell (posh)",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	code := 0

	// handle interrupt
	osInterrupt := make(chan os.Signal, 1)
	signal.Notify(osInterrupt, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	// handle defer
	defer func() {
		signal.Stop(osInterrupt)
		cancel()
		os.Exit(code)
	}()

	go func() {
		<-osInterrupt
		l.Debug("received interrupt")
		cancel()
	}()

	if err := rootCmd.ExecuteContext(ctx); errors.Is(err, context.Canceled) {
		l.Warn(err.Error())
	} else if err != nil {
		l.Error(err.Error())
		code = 1
	}
}

func init() {
	cobra.OnInitialize(initialize)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "disabled colors (default is false)")
	rootCmd.PersistentFlags().StringVar(&flagLevel, "level", "info", "set log level (default is warn)")
	rootCmd.PersistentFlags().StringVar(&flagConfig, "config", "", "config file (default is $HOME/.posh.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().Bool("no-validate", false, "Skip validation")
}

// initialize reads in config file and ENV variables if set.
func initialize() {
	// setup logger
	if value, err := log.NewPTerm(
		log.PTermWithDisableColor(flagNoColor),
		log.PTermWithLevel(log.GetLevel(flagLevel)),
	); err != nil {
		cobra.CheckErr(err)
	} else {
		l = value
	}

	// setup viper
	if flagConfig != "" {
		// Use config file from the flag.
		viper.SetConfigFile(flagConfig)
	} else {
		wd, err := os.Getwd()
		l.Must(err)
		viper.AddConfigPath(wd)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".posh")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	l.Must(viper.ReadInConfig())
	l.Debug("using config file:", viper.ConfigFileUsed())

	// validate version
	if v := viper.GetString("version"); v != "v1.0" {
		l.Must(fmt.Errorf("invalid config version: %s (v1.0)", v))
	}

	// setup env
	if value := os.Getenv(EnvProjectRoot); value != "" {
		// continue
	} else if value, err := os.Getwd(); err != nil {
		l.Must(errors.Wrap(err, "failed to retrieve project root"))
	} else if err := os.Setenv(EnvProjectRoot, value); err != nil {
		l.Must(errors.Wrap(err, "failed to set project root env"))
	}
	for key, value := range viper.GetStringMapString("env") {
		l.Must(os.Setenv(key, os.ExpandEnv(value)))
	}

	// setup manager
	var err error
	m, err = plugin.NewManager(l)
	l.Must(err)
}

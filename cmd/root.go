package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/foomo/posh/internal/env"
	intlog "github.com/foomo/posh/internal/log"
	"github.com/foomo/posh/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	l           log.Logger
	flagLevel   string
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
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "disabled colors (default is false)")
	rootCmd.PersistentFlags().StringVar(&flagLevel, "level", "info", "set log level (default is warn)")
}

// initialize reads in config file and ENV variables if set.
func initialize() {
	var err error

	// init logger
	l, err = intlog.Init(flagLevel, flagNoColor)
	cobra.CheckErr(err)

	// init env
	l.Must(env.Init())
}

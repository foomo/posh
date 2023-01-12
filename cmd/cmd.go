package cmd

import (
	"context"
	"os"
	"os/signal"

	intenv "github.com/foomo/posh/internal/env"
	intlog "github.com/foomo/posh/internal/log"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Init(provider plugin.Provider) {
	pluginProvider = provider
	cobra.OnInitialize(func() {
		l = intlog.Init(flagLevel, flagNoColor)
		l.Must(intenv.Init())
	})
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "disabled colors (default is false)")
	rootCmd.PersistentFlags().StringVar(&flagLevel, "level", "info", "set log level (default is warn)")
	rootCmd.AddCommand(
		configCmd,
		versionCmd,
	)

	if provider != nil {
		rootCmd.AddCommand(
			brewCmd,
			execCmd,
			promptCmd,
			requireCmd,
		)
		brewCmd.Flags().BoolVar(&brewCmdFlagDry, "dry", false, "don't execute scripts")
	}
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

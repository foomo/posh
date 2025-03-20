package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"

	cowsay "github.com/Code-Hex/Neo-cowsay/v2"
	"github.com/foomo/posh/internal/cmd"
	intenv "github.com/foomo/posh/internal/env"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Init(provider plugin.Provider) {
	pluginProvider = provider
	rootCmd = NewRoot()
	NewConfig(rootCmd)
	NewVersion(rootCmd)

	if provider != nil {
		NewBrew(rootCmd)
		NewExecute(rootCmd)
		NewPrompt(rootCmd)
		NewRequire(rootCmd)
	} else {
		NewInit(rootCmd)
	}

	cobra.OnInitialize(func() {
		if err := intenv.Init(); err != nil {
			panic(err)
		}
	})
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	code := 0
	l := cmd.NewLogger()

	// handle interrupt
	osInterrupt := make(chan os.Signal, 1)
	signal.Notify(osInterrupt, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	say := func(msg string) string {
		if say, err := cowsay.Say(msg, cowsay.BallonWidth(80)); err == nil {
			msg = say
		}
		return msg
	}

	// handle defer
	defer func() {
		if r := recover(); r != nil {
			l.Error(say("It's time to panic"))
			l.Error(fmt.Sprintf("%v", r))
			l.Error(string(debug.Stack()))
			code = 1
		}
		signal.Stop(osInterrupt)
		cancel()
		os.Exit(code)
	}()

	go func() {
		<-osInterrupt
		l.Trace("received interrupt")
		cancel()
	}()

	if err := rootCmd.ExecuteContext(ctx); errors.Is(err, context.Canceled) {
		l.Warn(err.Error())
	} else if err != nil {
		l.Error(say(strings.Split(errors.Cause(err).Error(), ":")[0]))
		l.Error(err.Error())
		code = 1
	}
}

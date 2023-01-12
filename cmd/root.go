package cmd

import (
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/spf13/cobra"
)

var (
	l              log.Logger
	flagLevel      string
	flagNoColor    bool
	pluginProvider plugin.Provider
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "posh",
	Short: "Project Oriented Shell (posh)",
}

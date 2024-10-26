package cmd

import (
	intconfig "github.com/foomo/posh/internal/config"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd        *cobra.Command
	pluginProvider plugin.Provider
)

// NewRoot represents the base command when called without any subcommands
func NewRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "posh",
		Short: "Project Oriented Shell (posh)",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return intconfig.Dotenv()
		},
	}

	cmd.PersistentFlags().Bool("no-color", false, "disabled colors (default: false)")
	_ = viper.BindPFlag("no-color", cmd.PersistentFlags().Lookup("no-color"))

	cmd.PersistentFlags().String("level", "info", "set log level (default: info)")
	_ = viper.BindPFlag("level", cmd.PersistentFlags().Lookup("level"))

	return cmd
}

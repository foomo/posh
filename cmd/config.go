package cmd

import (
	intconfig "github.com/foomo/posh/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:           "config",
	Short:         "Print loaded configuration",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := intconfig.Load(l); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Debug()
		return nil
	},
}

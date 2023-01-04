package cmd

import (
	intconfig "github.com/foomo/posh/internal/config"
	intplugin "github.com/foomo/posh/internal/plugin"
	"github.com/foomo/posh/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dependenciesCmd represents the dependencies command
var dependenciesCmd = &cobra.Command{
	Use:           "dependencies",
	Short:         "Run dependency validations",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := intconfig.Load(l, flagConfig); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Dependencies
		if err := viper.UnmarshalKey("dependencies", &cfg); err != nil {
			return err
		}

		plg, err := intplugin.Load(cmd.Context(), l)
		if err != nil {
			return err
		}

		return plg.Dependencies(cmd.Context(), cfg)
	},
}

func init() {
	rootCmd.AddCommand(dependenciesCmd)
}

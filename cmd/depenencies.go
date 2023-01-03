package cmd

import (
	"github.com/foomo/posh/internal/util"
	"github.com/foomo/posh/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dependenciesCmd represents the dependencies command
var dependenciesCmd = &cobra.Command{
	Use:           "dependencies",
	Short:         "Run dependency validations",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Dependencies
		if err := viper.UnmarshalKey("dependencies", &cfg); err != nil {
			return err
		}

		plg, err := util.LoadPlugin(cmd.Context(), m)
		if err != nil {
			return err
		}

		return plg.Dependencies(cmd.Context(), cfg)
	},
}

func init() {
	rootCmd.AddCommand(dependenciesCmd)
}

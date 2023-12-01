package cmd

import (
	intconfig "github.com/foomo/posh/internal/config"
	"github.com/foomo/posh/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// requireCmd represents the require command
var requireCmd = &cobra.Command{
	Use:           "require",
	Short:         "Validate configured requirements",
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := intconfig.Load(l); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Require
		if err := viper.UnmarshalKey("require", &cfg); err != nil {
			return err
		}

		plg, err := pluginProvider(l)
		if err != nil {
			return err
		}

		return plg.Require(cmd.Context(), cfg)
	},
}

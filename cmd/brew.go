package cmd

import (
	intconfig "github.com/foomo/posh/internal/config"
	"github.com/foomo/posh/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	brewCmdFlagDry bool
)

// brewCmd represents the dependencies command
var brewCmd = &cobra.Command{
	Use:           "brew",
	Short:         "Check and install required packages.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := intconfig.Load(l); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Ownbrew
		if err := viper.UnmarshalKey("ownbrew", &cfg); err != nil {
			return err
		}
		cfg.Dry = brewCmdFlagDry

		plg, err := pluginProvider(l)
		if err != nil {
			return err
		}

		return plg.Brew(cmd.Context(), cfg)
	},
}

package cmd

import (
	ownbrewconfig "github.com/foomo/ownbrew/pkg/config"
	"github.com/foomo/posh/internal/config"
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
		if err := config.Load(l); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg ownbrewconfig.Config
		if err := viper.UnmarshalKey("ownbrew", &cfg); err != nil {
			return err
		}

		plg, err := pluginProvider(l)
		if err != nil {
			return err
		}

		return plg.Brew(cmd.Context(), cfg, brewCmdFlagDry)
	},
}

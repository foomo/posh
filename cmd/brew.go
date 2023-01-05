package cmd

import (
	intconfig "github.com/foomo/posh/internal/config"
	intplugin "github.com/foomo/posh/internal/plugin"
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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := intconfig.Load(l, flagConfig); err != nil {
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

		plg, err := intplugin.Load(cmd.Context(), l)
		if err != nil {
			return err
		}

		return plg.Brew(cmd.Context(), cfg)
	},
}

func init() {
	rootCmd.AddCommand(brewCmd)
	brewCmd.Flags().BoolVar(&brewCmdFlagDry, "dry", false, "don't execute scripts")
}

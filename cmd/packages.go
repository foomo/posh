package cmd

import (
	"github.com/foomo/posh/internal/util"
	"github.com/foomo/posh/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// packagesCmd represents the dependencies command
var packagesCmd = &cobra.Command{
	Use:           "packages",
	Short:         "Check and install required packages.",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg []config.Package
		if err := viper.UnmarshalKey("packages", &cfg); err != nil {
			return err
		}

		plg, err := util.LoadPlugin(cmd.Context(), m)
		if err != nil {
			return err
		}

		return plg.Packages(cmd.Context(), cfg)
	},
}

func init() {
	rootCmd.AddCommand(packagesCmd)
}

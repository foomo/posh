package cmd

import (
	"github.com/foomo/posh/internal/util"
	"github.com/foomo/posh/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:           "prompt",
	Short:         "Start the interactive Project Oriented Shell",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Prompt
		if err := viper.UnmarshalKey("prompt", &cfg); err != nil {
			return err
		}

		plg, err := util.LoadPlugin(cmd.Context(), m)
		if err != nil {
			return err
		}

		return plg.Prompt(cmd.Context(), cfg)
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
}

package cmd

import (
	"context"

	intconfig "github.com/foomo/posh/internal/config"
	"github.com/foomo/posh/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:           "prompt",
	Short:         "Start the interactive Project Oriented Shell",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := intconfig.Load(l); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Prompt
		if err := viper.UnmarshalKey("prompt", &cfg); err != nil {
			return err
		}

		plg, err := pluginProvider(l)
		if err != nil {
			return err
		}

		return plg.Prompt(context.TODO(), cfg)
	},
}

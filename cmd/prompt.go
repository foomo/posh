package cmd

import (
	intcmd "github.com/foomo/posh/internal/cmd"
	intconfig "github.com/foomo/posh/internal/config"
	"github.com/foomo/posh/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewPrompt represents the prompt command
func NewPrompt(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:           "prompt",
		Short:         "Start the interactive Project Oriented Shell",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			l := intcmd.NewLogger()
			if err := intconfig.Load(l); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			l := intcmd.NewLogger()
			var cfg config.Prompt
			if err := viper.UnmarshalKey("prompt", &cfg); err != nil {
				return err
			}

			plg, err := pluginProvider(l)
			if err != nil {
				return err
			}

			return plg.Prompt(cmd.Context(), cfg)
		},
	}

	root.AddCommand(cmd)
}

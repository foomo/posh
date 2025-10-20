package cmd

import (
	intcmd "github.com/foomo/posh/internal/cmd"
	intconfig "github.com/foomo/posh/internal/config"
	"github.com/foomo/posh/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewRequire represents the require command
func NewRequire(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:           "require",
		Short:         "Validate configured requirements",
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
			var cfg config.Require

			l := intcmd.NewLogger()

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

	root.AddCommand(cmd)
}

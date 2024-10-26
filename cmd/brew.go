package cmd

import (
	ownbrewconfig "github.com/foomo/ownbrew/pkg/config"
	intcmd "github.com/foomo/posh/internal/cmd"
	"github.com/foomo/posh/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewBrew represents the dependencies command
func NewBrew(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:           "brew",
		Short:         "Check and install required packages.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			l := intcmd.NewLogger()
			if err := config.Load(l); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			l := intcmd.NewLogger()
			var cfg ownbrewconfig.Config
			if err := viper.UnmarshalKey("ownbrew", &cfg); err != nil {
				return err
			}

			plg, err := pluginProvider(l)
			if err != nil {
				return err
			}

			dry, err := cmd.Flags().GetBool("dry")
			if err != nil {
				return err
			}

			tags, err := cmd.Flags().GetStringSlice("tags")
			if err != nil {
				return err
			}

			return plg.Brew(cmd.Context(), cfg, tags, dry)
		},
	}

	cmd.Flags().Bool("dry", false, "print out the taps that will be installed")
	_ = viper.BindPFlag("dry", cmd.Flags().Lookup("dry"))

	cmd.Flags().StringSlice("tags", nil, "filter by tags (e.g. ci,-test)")
	_ = viper.BindPFlag("tags", cmd.Flags().Lookup("tags"))

	root.AddCommand(cmd)
}

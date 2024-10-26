package cmd

import (
	intcmd "github.com/foomo/posh/internal/cmd"
	intconfig "github.com/foomo/posh/internal/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewExecute represents the exec command
func NewExecute(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:                "execute",
		Short:              "Execute a single Project Oriented Shell command",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			l := intcmd.NewLogger()
			if err := intconfig.Load(l); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			l := intcmd.NewLogger()
			if len(args) == 0 {
				return errors.New("missing [cmd] argument")
			}

			plg, err := pluginProvider(l)
			if err != nil {
				return err
			}

			return plg.Execute(cmd.Context(), args)
		},
	}

	root.AddCommand(cmd)
}

package cmd

import (
	intcmd "github.com/foomo/posh/internal/cmd"
	intconfig "github.com/foomo/posh/internal/config"
	"github.com/foomo/posh/pkg/agent"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/foomo/posh/pkg/readline"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewExecute represents the exec command
func NewExecute(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:                "execute",
		Short:              "Execute a single Project Oriented Shell command",
		Aliases:            []string{"x"},
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
			// DisableFlagParsing means cobra never populates this command's
			// flags, so posh's own are scanned off the front by hand.
			args = readline.ExtractFlags(args, map[string]func(){
				agent.Flag: func() { agent.SetFlag(true) },
			})

			// Logger must be built after ExtractFlags, which may have set
			// agent mode via a leading --agent.
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
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			l := intcmd.NewLogger()
			if err := intconfig.Load(l); err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			plg, err := pluginProvider(l)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			completer, ok := plg.(plugin.Completer)
			if !ok {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			return completer.Complete(cmd.Context(), args, toComplete), cobra.ShellCompDirectiveNoFileComp
		},
	}

	root.AddCommand(cmd)
}

package cmd

import (
	intcmd "github.com/foomo/posh/internal/cmd"
	intconfig "github.com/foomo/posh/internal/config"
	"github.com/foomo/posh/pkg/mcp"
	"github.com/spf13/cobra"
)

// NewMCP represents the mcp command
func NewMCP(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:           "mcp",
		Short:         "Serve an MCP server over stdio for this project's posh shell",
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

			return mcp.New(l, pluginProvider).Run(cmd.Context())
		},
	}

	root.AddCommand(cmd)
}

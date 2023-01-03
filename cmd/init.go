package cmd

import (
	"os"
	"path/filepath"

	"github.com/foomo/posh/embed"
	"github.com/foomo/posh/pkg/scaffold"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a Project Oriented Shell",
	Long: `Initialize (posh init) will create a new Project Oriented Shell with the appropriate structure.

Posh init must be run inside of a go module (please run "go mod init <MODNAME> first)"`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		if len(args) > 0 && args[0] != "." {
			wd = filepath.Join(wd, args[0])
		}

		fs, err := embed.Scaffold("init")
		if err != nil {
			return err
		}

		sc, err := scaffold.New(
			scaffold.WithLogger(l),
		)
		if err != nil {
			return err
		}

		return sc.Render(fs, wd, nil)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

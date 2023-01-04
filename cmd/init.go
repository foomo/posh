package cmd

import (
	"os"
	"path"

	"github.com/foomo/posh/embed"
	"github.com/foomo/posh/internal/util"
	"github.com/foomo/posh/pkg/env"
	"github.com/foomo/posh/pkg/scaffold"
	"github.com/spf13/cobra"
)

var (
	initCmdFlagDry bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Project Oriented Shell",
	Long: `Initialize (posh init) will create a new Project Oriented Shell with the appropriate structure.

Posh init must be run inside of a go module (please run "go mod init <MODNAME> first)"`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		data := map[string]interface{}{}

		// define module
		if gitRemoteURL, err := util.GitRemoteURL(); err == nil {
			data["module"] = gitRemoteURL
		} else {
			data["module"] = path.Base(os.Getenv(env.ProjectRoot))
		}

		fs, err := embed.Scaffold("init")
		if err != nil {
			return err
		}

		sc, err := scaffold.New(
			scaffold.WithLogger(l),
			scaffold.WithDry(initCmdFlagDry),
			scaffold.WithDirectories(scaffold.Directory{
				Source: fs,
				Target: os.Getenv(env.ProjectRoot),
				Data:   data,
			}),
		)
		if err != nil {
			return err
		}

		return sc.Render(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&initCmdFlagDry, "dry", false, "don't render files")
}

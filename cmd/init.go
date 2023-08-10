package cmd

import (
	"path"

	"github.com/foomo/posh/embed"
	scaffold2 "github.com/foomo/posh/integration/scaffold"
	"github.com/foomo/posh/internal/util/git"
	"github.com/foomo/posh/pkg/env"
	"github.com/spf13/cobra"
)

var (
	initCmdFlagDry      bool
	initCmdFlagOverride bool
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
		if value, err := git.OriginURL(); err != nil {
			l.Debug("failed to retrieve git origin url:", err.Error())
			data["module"] = path.Base(env.ProjectRoot())
		} else {
			data["module"] = value
		}

		fs, err := embed.Scaffold("init")
		if err != nil {
			return err
		}

		sc, err := scaffold2.New(
			l,
			scaffold2.WithDry(initCmdFlagDry),
			scaffold2.WithOverride(initCmdFlagOverride),
			scaffold2.WithDirectories(scaffold2.Directory{
				Source: fs,
				Target: env.ProjectRoot(),
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
	initCmd.Flags().BoolVar(&initCmdFlagOverride, "override", false, "override existing files")
}

package cmd

import (
	"path"

	"github.com/foomo/posh/embed"
	scaffold2 "github.com/foomo/posh/integration/scaffold"
	intcmd "github.com/foomo/posh/internal/cmd"
	"github.com/foomo/posh/internal/util/git"
	"github.com/foomo/posh/pkg/env"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewInit represents the init command
func NewInit(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a Project Oriented Shell",
		Long: `Initialize (posh init) will create a new Project Oriented Shell with the appropriate structure.

Posh init must be run inside of a go module (please run "go mod init <MODNAME> first)"`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			l := intcmd.NewLogger()
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

			dry, err := cmd.Flags().GetBool("dry")
			if err != nil {
				return err
			}

			override, err := cmd.Flags().GetBool("override")
			if err != nil {
				return err
			}

			sc, err := scaffold2.New(
				l,
				scaffold2.WithDry(dry),
				scaffold2.WithOverride(override),
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

	cmd.Flags().Bool("dry", false, "don't render files")
	_ = viper.BindPFlag("dry", cmd.Flags().Lookup("dry"))

	cmd.Flags().Bool("override", false, "override existing files")
	_ = viper.BindPFlag("override", cmd.Flags().Lookup("override"))

	root.AddCommand(cmd)
}

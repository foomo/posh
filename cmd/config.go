package cmd

import (
	"fmt"

	"github.com/foomo/ownbrew/pkg/util"
	intcmd "github.com/foomo/posh/internal/cmd"
	intconfig "github.com/foomo/posh/internal/config"
	"github.com/foomo/posh/pkg/agent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// NewConfig represents the config command
func NewConfig(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:           "config",
		Short:         "Print loaded configuration",
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
			// The settings map is already what an agent wants, so it is encoded
			// as-is rather than through a payload struct: the config shape is
			// open ended and defined by each project, not by posh.
			return agent.Render(
				func() any { return viper.AllSettings() },
				func() error {
					out, err := yaml.Marshal(viper.AllSettings())
					if err != nil {
						return err
					}

					// Highlight colors via chroma, which --no-color does not
					// reach: the flag only toggles a pterm global, and this
					// command writes to stdout rather than through the logger.
					if viper.GetBool("no-color") {
						fmt.Println(string(out))
					} else {
						fmt.Println(util.Highlight(string(out), "yaml"))
					}

					return nil
				},
			)
		},
	}

	root.AddCommand(cmd)
}

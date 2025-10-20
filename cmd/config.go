package cmd

import (
	"fmt"

	"github.com/foomo/ownbrew/pkg/util"
	intcmd "github.com/foomo/posh/internal/cmd"
	intconfig "github.com/foomo/posh/internal/config"
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
			out, err := yaml.Marshal(viper.AllSettings())
			if err != nil {
				return err
			}

			fmt.Println(util.Highlight(string(out), "yaml"))

			return nil
		},
	}

	root.AddCommand(cmd)
}

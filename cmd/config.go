package cmd

import (
	"fmt"

	intconfig "github.com/foomo/posh/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:           "config",
	Short:         "Print loaded configuration",
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
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
		fmt.Println(string(out))
		return nil
	},
}

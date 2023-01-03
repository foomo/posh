package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "develop"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  `If unsure which version of the CLI you are using, you can use this command to print the version of the CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

package cmd

import (
	intversion "github.com/foomo/posh/internal/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  `If unsure which version of the CLI you are using, you can use this command to print the version of the CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		l.Debugf("%s (%s)", intversion.CommitHash, intversion.BuildTimestamp)
		l.Print(intversion.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

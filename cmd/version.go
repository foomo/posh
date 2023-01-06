package cmd

import (
	"strconv"
	"time"

	intversion "github.com/foomo/posh/internal/version"
	"github.com/foomo/posh/pkg/log"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  `If unsure which version of the CLI you are using, you can use this command to print the version of the CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		buildTime := intversion.BuildTimestamp
		if value, err := strconv.ParseInt(intversion.BuildTimestamp, 10, 64); err == nil {
			buildTime = time.Unix(value, 0).String()
		}
		if l.IsLevel(log.LevelDebug) {
			l.Printf("v%s, Commit: %s, BuildTime: %s", intversion.Version, intversion.CommitHash, buildTime)
		} else {
			l.Printf("v%s", intversion.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

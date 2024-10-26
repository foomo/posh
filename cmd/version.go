package cmd

import (
	"strconv"
	"time"

	intcmd "github.com/foomo/posh/internal/cmd"
	intversion "github.com/foomo/posh/internal/version"
	"github.com/foomo/posh/pkg/log"
	"github.com/spf13/cobra"
)

// NewVersion represents the version command
func NewVersion(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Long:  `If unsure which version of the CLI you are using, you can use this command to print the version of the CLI.`,
		Run: func(cmd *cobra.Command, args []string) {
			l := intcmd.NewLogger()
			buildTime := intversion.BuildTimestamp
			if value, err := strconv.ParseInt(intversion.BuildTimestamp, 10, 64); err == nil {
				buildTime = time.Unix(value, 0).String()
			}
			if l.IsLevel(log.LevelDebug) {
				l.Printf("Version: %s\nCommit: %s\nBuildTime: %s", intversion.Version, intversion.CommitHash, buildTime)
			} else {
				l.Printf("%s", intversion.Version)
			}
		},
	}

	root.AddCommand(cmd)
}

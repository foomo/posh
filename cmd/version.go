package cmd

import (
	"strconv"
	"time"

	intcmd "github.com/foomo/posh/internal/cmd"
	intversion "github.com/foomo/posh/internal/version"
	"github.com/foomo/posh/pkg/agent"
	"github.com/foomo/posh/pkg/command"
	"github.com/foomo/posh/pkg/log"
	"github.com/spf13/cobra"
)

// NewVersion represents the version command
func NewVersion(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:           "version",
		Short:         "Print the version",
		Long:          `If unsure which version of the CLI you are using, you can use this command to print the version of the CLI.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			l := intcmd.NewLogger()

			buildTime := intversion.BuildTimestamp
			if value, err := strconv.ParseInt(intversion.BuildTimestamp, 10, 64); err == nil {
				buildTime = time.Unix(value, 0).String()
			}

			// The debug gate governs both forms, so the agent payload carries
			// exactly the fields the human output prints.
			debug := l.IsLevel(log.LevelDebug)

			return agent.Render(
				func() any {
					ret := command.Version{Version: intversion.Version}
					if debug {
						ret.Commit = intversion.CommitHash
						ret.BuildTime = buildTime
					}

					return ret
				},
				func() error {
					if debug {
						l.Printf("Version: %s\nCommit: %s\nBuildTime: %s", intversion.Version, intversion.CommitHash, buildTime)
					} else {
						l.Printf("%s", intversion.Version)
					}

					return nil
				},
			)
		},
	}

	root.AddCommand(cmd)
}

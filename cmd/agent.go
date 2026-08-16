package cmd

import (
	"context"

	intcmd "github.com/foomo/posh/internal/cmd"
	intconfig "github.com/foomo/posh/internal/config"
	"github.com/foomo/posh/pkg/agent"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/spf13/cobra"
)

// NewAgent represents the agent command group: tooling for AI coding agents
// driving posh, as opposed to the per-invocation `--agent` root flag.
func NewAgent(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Tooling for AI coding agents driving this project's posh shell",
	}

	newAgentCatalog(cmd)
	newAgentSkill(cmd)

	root.AddCommand(cmd)
}

// newAgentCatalog represents the `agent catalog` command
func newAgentCatalog(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:           "catalog",
		Short:         "Print this project's command catalog as JSON",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			l := intcmd.NewLogger()
			return intconfig.Load(l)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			commands, err := loadCatalog(cmd.Context())
			if err != nil {
				return err
			}

			return agent.Encode(plugin.NewCatalog(commands))
		},
	}

	root.AddCommand(cmd)
}

// newAgentSkill represents the `agent skill` command group
func newAgentSkill(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Generate a Claude Code skill describing this project's posh shell",
	}

	install := &cobra.Command{
		Use:           "install [path]",
		Short:         "Write the generated skill to disk",
		Args:          cobra.MaximumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			l := intcmd.NewLogger()
			return intconfig.Load(l)
		},
		RunE: runAgentSkillInstall,
	}

	// get renders the same skill install would write, but to stdout - so an
	// agent can read it without touching the working tree.
	get := &cobra.Command{
		Use:           "get",
		Short:         "Print the generated skill to stdout",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE:       install.PreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			meta, commands, err := loadSkill(cmd.Context())
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(plugin.RenderSkill(meta, commands)))

			return err
		},
	}

	uninstall := &cobra.Command{
		Use:           "uninstall [path]",
		Short:         "Remove the generated skill from disk",
		Args:          cobra.MaximumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var path string
			if len(args) > 0 {
				path = args[0]
			}

			return plugin.RemoveSkill(path)
		},
	}

	// update is install by another name: both (re-)generate the skill from
	// this project's live command catalog, so re-running after adding or
	// changing providers always produces an up-to-date file.
	update := &cobra.Command{
		Use:           "update [path]",
		Short:         "Regenerate the skill file after this project's commands changed",
		Args:          cobra.MaximumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE:       install.PreRunE,
		RunE:          runAgentSkillInstall,
	}

	cmd.AddCommand(get, install, update, uninstall)
	root.AddCommand(cmd)
}

func runAgentSkillInstall(cmd *cobra.Command, args []string) error {
	meta, commands, err := loadSkill(cmd.Context())
	if err != nil {
		return err
	}

	var path string
	if len(args) > 0 {
		path = args[0]
	}

	return plugin.WriteSkill(path, meta, commands)
}

// loadCatalog resolves this project's plugin and asks it for the command
// catalog.
//
// Only the resolution stays here: it goes through pluginProvider - package
// state set by Init - and internal/cmd, neither of which pkg/ may import.
func loadCatalog(ctx context.Context) ([]plugin.CommandInfo, error) {
	plg, err := pluginProvider(intcmd.NewLogger())
	if err != nil {
		return nil, err
	}

	return plugin.List(ctx, plg)
}

// loadSkill resolves this project's plugin and asks it for everything the
// generated skill needs: the frontmatter and the command catalog with per
// command prose.
//
// It is kept apart from loadCatalog so `agent catalog` keeps emitting the
// unchanged catalog shape.
func loadSkill(ctx context.Context) (plugin.SkillMetadata, []plugin.SkillCommand, error) {
	plg, err := pluginProvider(intcmd.NewLogger())
	if err != nil {
		return plugin.SkillMetadata{}, nil, err
	}

	commands, err := plugin.ListSkill(ctx, plg)
	if err != nil {
		return plugin.SkillMetadata{}, nil, err
	}

	return plugin.SkillMetadataOf(ctx, plg), commands, nil
}

package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

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
		Short:         "Write the generated skills to disk",
		Args:          cobra.MaximumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			l := intcmd.NewLogger()
			return intconfig.Load(l)
		},
		RunE: runAgentSkillInstall,
	}

	// get renders the same skills install would write, but to stdout - so an
	// agent can read them without touching the working tree.
	//
	// Each one is preceded by its path, since there are now several: without it
	// the frontmatter of the next skill is the only clue a new file started.
	get := &cobra.Command{
		Use:           "get",
		Short:         "Print the generated skills to stdout",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE:       install.PreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			meta, commands, err := loadSkill(cmd.Context())
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()

			skills := []struct{ path, content string }{{
				path:    filepath.Join(plugin.RootSkillName, "SKILL.md"),
				content: plugin.RenderRootSkill(meta, commands),
			}}

			for _, c := range commands {
				if !c.Describes() {
					continue
				}

				skills = append(skills, struct{ path, content string }{
					path:    filepath.Join(c.SkillName(), "SKILL.md"),
					content: plugin.RenderCommandSkill(c),
				})
			}

			for _, skill := range skills {
				if _, err := fmt.Fprintf(out, "==> %s\n\n%s\n", skill.path, skill.content); err != nil {
					return err
				}
			}

			return nil
		},
	}

	uninstall := &cobra.Command{
		Use:           "uninstall [path]",
		Short:         "Remove the generated skills from disk",
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

	// update is install by another name: both (re-)generate the skills from
	// this project's live command catalog, so re-running after adding or
	// changing providers always produces up-to-date files.
	update := &cobra.Command{
		Use:           "update [path]",
		Short:         "Regenerate the skill files after this project's commands changed",
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

	written, err := plugin.WriteSkill(path, meta, commands)
	if err != nil {
		return err
	}

	l := intcmd.NewLogger()

	l.Successf("wrote %d skills to %s", len(written), plugin.SkillsPath(path))

	// A command without its own description gets one derived from its one-line
	// description, which says what it is rather than when to use it - and the
	// description is the only thing an agent runtime matches against when
	// deciding to load a skill. Report it: the skill is generated either way, so
	// the gap is otherwise invisible.
	if fallback := plugin.FallbackSkillDescriptions(commands); len(fallback) > 0 {
		l.Warnf(
			"%d commands have no skill description and fell back to a generated one, so their skills may never load: %s",
			len(fallback), strings.Join(fallback, ", "),
		)
		l.Info("implement command.SkillMetadataer on these commands to name the conditions that should reach for them")
	}

	return nil
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

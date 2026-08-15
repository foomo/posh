package plugin_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	ownbrewconfig "github.com/foomo/ownbrew/pkg/config"
	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestRenderSkill(t *testing.T) {
	out := plugin.RenderSkill(plugin.SkillMetadata{}, []plugin.SkillCommand{
		{
			CommandInfo: plugin.CommandInfo{
				FullPath:    "kubectl",
				Description: "manage kubernetes clusters",
				Subcommands: []plugin.CommandInfo{
					{
						FullPath:    "kubectl <cluster>",
						Description: "target cluster",
						Dynamic:     true,
						Subcommands: []plugin.CommandInfo{
							{
								FullPath:    "kubectl <cluster> apply",
								Description: "apply manifests",
								Arguments: []plugin.ArgInfo{
									{Name: "Path", Description: "manifest directory", Repeat: true},
								},
								Flags: []plugin.FlagInfo{
									{Name: "dry-run", Type: "bool", Description: "simulate the apply"},
								},
							},
						},
					},
				},
			},
		},
		{CommandInfo: plugin.CommandInfo{FullPath: "exit", Description: "exit shell"}},
	})

	assert.Contains(t, out, "name: posh")
	assert.Contains(t, out, "### `kubectl`")
	assert.Contains(t, out, "manage kubernetes clusters")

	// exit contributes neither prose nor structure, so it is omitted.
	assert.NotContains(t, out, "### `exit`")

	// nested descendants are rendered, deepest carrying args and flags
	assert.Contains(t, out, "- `kubectl <cluster>` - target cluster")
	assert.Contains(t, out, "  - `kubectl <cluster> apply [Path]...` - apply manifests")
	assert.Contains(t, out, "`[Path]` (repeatable) manifest directory")
	assert.Contains(t, out, "`--dry-run` (bool) simulate the apply")
}

// TestRenderSkill_UsageHeading covers the section heading: it separates a
// command's structure from its description, and is omitted when there is no
// structure to head.
func TestRenderSkill_UsageHeading(t *testing.T) {
	out := plugin.RenderSkill(plugin.SkillMetadata{}, []plugin.SkillCommand{
		{
			CommandInfo: plugin.CommandInfo{
				FullPath:    "env",
				Description: "Manage internal environment variables",
				Subcommands: []plugin.CommandInfo{
					{FullPath: "env list", Description: "List all environment variables"},
				},
			},
			Skill: "#### Notes\n\nProcess scoped.",
		},
	})

	// The heading sits between the description and the list, and the command's
	// own prose keeps its heading after it.
	assert.Regexp(t,
		"(?s)Manage internal environment variables\n+#### Usage\n+- `env list`.*#### Notes",
		out)

	// A command whose only contribution is prose has no structure to head.
	out = plugin.RenderSkill(plugin.SkillMetadata{}, []plugin.SkillCommand{
		{
			CommandInfo: plugin.CommandInfo{FullPath: "check", Description: "run checks"},
			Skill:       "#### Notes\n\nRead the results.",
		},
	})

	assert.NotContains(t, out, "#### Usage")
	assert.Contains(t, out, "#### Notes")
}

// TestRenderSkill_OmitsBareLeaves covers the filter: a command carrying nothing
// but a name and description tells an agent no more than the catalog already
// does, so it is left out - unless it contributes prose.
func TestRenderSkill_OmitsBareLeaves(t *testing.T) {
	out := plugin.RenderSkill(plugin.SkillMetadata{}, []plugin.SkillCommand{
		{CommandInfo: plugin.CommandInfo{FullPath: "bare", Description: "nothing to add"}},
		{
			CommandInfo: plugin.CommandInfo{FullPath: "prose", Description: "leaf with notes"},
			Skill:       "#### Notes\n\nWorth saying.",
		},
		{
			CommandInfo: plugin.CommandInfo{
				FullPath:    "args",
				Description: "leaf with an argument",
				Arguments:   []plugin.ArgInfo{{Name: "Key"}},
			},
		},
		{
			CommandInfo: plugin.CommandInfo{
				FullPath:    "flags",
				Description: "leaf with a flag",
				Flags:       []plugin.FlagInfo{{Name: "dry-run", Type: "bool"}},
			},
		},
	})

	assert.NotContains(t, out, "### `bare`")

	// Any one of prose, arguments or flags is enough to keep a command.
	assert.Contains(t, out, "### `prose`")
	assert.Contains(t, out, "### `args")
	assert.Contains(t, out, "### `flags`")

	// Whitespace-only prose does not count as a contribution.
	out = plugin.RenderSkill(plugin.SkillMetadata{}, []plugin.SkillCommand{
		{CommandInfo: plugin.CommandInfo{FullPath: "blank"}, Skill: "  \n\t\n"},
	})
	assert.NotContains(t, out, "### `blank`")
}

func TestRenderSkill_TopLevelArgsAndFlags(t *testing.T) {
	out := plugin.RenderSkill(plugin.SkillMetadata{}, []plugin.SkillCommand{
		{
			CommandInfo: plugin.CommandInfo{
				FullPath:    "welcome",
				Description: "print a welcome message",
				Arguments: []plugin.ArgInfo{
					{Name: "Name", Description: "who to greet", Optional: true},
				},
				Flags: []plugin.FlagInfo{
					{Name: "loud", Type: "bool", Description: "shout it"},
				},
			},
		},
	})

	// The usage line is rebuilt from the structured arguments.
	assert.Contains(t, out, "### `welcome <Name>`")
	assert.Contains(t, out, "- `<Name>` who to greet")
	assert.Contains(t, out, "- `--loud` (bool) shout it")
}

func TestRenderSkill_Empty(t *testing.T) {
	out := plugin.RenderSkill(plugin.SkillMetadata{}, nil)

	assert.Contains(t, out, "## Commands")
	assert.NotContains(t, out, "###")
}

// TestRenderSkill_Setup covers posh's own subcommands, which are documented from
// a fixed string rather than the catalog - so they render even for a project
// whose plugin registers no commands at all.
func TestRenderSkill_Setup(t *testing.T) {
	out := plugin.RenderSkill(plugin.SkillMetadata{}, nil)

	assert.Contains(t, out, "## Setup")
	assert.Contains(t, out, "`.posh.yaml`")
	assert.Contains(t, out, "`posh require`")
	assert.Contains(t, out, "`posh brew`")
	assert.Contains(t, out, "`posh prompt`")
}

// TestRenderSkill_Skill covers the command.Skiller contribution: the markdown is
// appended verbatim under the command's own heading.
func TestRenderSkill_Skill(t *testing.T) {
	out := plugin.RenderSkill(plugin.SkillMetadata{}, []plugin.SkillCommand{
		{
			CommandInfo: plugin.CommandInfo{FullPath: "cache", Description: "manage caches"},
			// The trailing newline must not stack onto the one the renderer
			// adds itself.
			Skill: "#### Configuration\n\n```yaml\ncache:\n  ttl: 5m\n```\n",
		},
		{
			CommandInfo: plugin.CommandInfo{
				FullPath:    "env",
				Description: "manage env vars",
				Arguments:   []plugin.ArgInfo{{Name: "Key"}},
			},
		},
	})

	assert.Contains(t, out, "#### Configuration")
	assert.Contains(t, out, "  ttl: 5m")
	assert.NotContains(t, out, "```\n\n\n", "a trailing newline must not stack blank lines")

	// A command contributing structure but no prose is unaffected.
	assert.Contains(t, out, "### `env [Key]`")
}

// TestRenderSkill_Frontmatter covers the SkillMetadata overrides, including the
// YAML-hostile description that hand-built frontmatter would corrupt.
func TestRenderSkill_Frontmatter(t *testing.T) {
	out := plugin.RenderSkill(plugin.SkillMetadata{
		Name:         "acme",
		Description:  "Drive acme: deploy, seed the DB #fast",
		AllowedTools: []string{"Bash(posh execute:*)"},
	}, nil)

	front, _, ok := strings.Cut(strings.TrimPrefix(out, "---\n"), "---\n")
	require.True(t, ok, "the frontmatter must be delimited")

	var actual plugin.SkillMetadata
	require.NoError(t, yaml.Unmarshal([]byte(front), &actual), "the frontmatter must be valid YAML")

	assert.Equal(t, "acme", actual.Name)
	assert.Equal(t, "Drive acme: deploy, seed the DB #fast", actual.Description)
	assert.Equal(t, []string{"Bash(posh execute:*)"}, actual.AllowedTools)

	// The hyphen form is what Claude Code reads; allowed_tools is not accepted.
	assert.Contains(t, front, "allowed-tools:")
}

// TestRenderSkill_FrontmatterDefaults covers the zero value reproducing posh's
// own frontmatter, with no empty allowed-tools key.
func TestRenderSkill_FrontmatterDefaults(t *testing.T) {
	out := plugin.RenderSkill(plugin.SkillMetadata{}, nil)

	assert.Contains(t, out, "name: posh")
	assert.Contains(t, out, "Drive this project's posh shell")
	assert.NotContains(t, out, "allowed-tools")
}

func TestWriteSkill(t *testing.T) {
	// A nested path exercises the parent-directory creation.
	path := filepath.Join(t.TempDir(), "nested", "dir", "SKILL.md")

	require.NoError(t, plugin.WriteSkill(path, plugin.SkillMetadata{}, []plugin.SkillCommand{
		{
			CommandInfo: plugin.CommandInfo{FullPath: "welcome", Description: "print a welcome message"},
			Skill:       "#### Notes\n\nNo side effects.",
		},
	}))

	b, err := os.ReadFile(path)
	require.NoError(t, err)

	assert.Contains(t, string(b), "name: posh")
	assert.Contains(t, string(b), "### `welcome`")
}

func TestWriteSkill_DefaultPath(t *testing.T) {
	// DefaultSkillPath is relative, so write it inside a temp working dir.
	dir := t.TempDir()
	t.Chdir(dir)

	require.NoError(t, plugin.WriteSkill("", plugin.SkillMetadata{}, nil))

	_, err := os.Stat(filepath.Join(dir, plugin.DefaultSkillPath))
	assert.NoError(t, err, "an empty path must fall back to DefaultSkillPath")
}

func TestRemoveSkill(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	require.NoError(t, plugin.WriteSkill("", plugin.SkillMetadata{}, nil))
	require.NoError(t, plugin.RemoveSkill(""))

	_, err := os.Stat(filepath.Join(dir, plugin.DefaultSkillPath))
	assert.True(t, os.IsNotExist(err), "the skill file must be gone")

	assert.NoError(t, plugin.RemoveSkill(""), "removing a missing file must be a no-op")
}

type stubLister struct{ commands []plugin.CommandInfo }

func (s stubLister) Prompt(ctx context.Context, cfg config.Prompt) error   { return nil }
func (s stubLister) Execute(ctx context.Context, args []string) error      { return nil }
func (s stubLister) Require(ctx context.Context, cfg config.Require) error { return nil }
func (s stubLister) Brew(ctx context.Context, cfg ownbrewconfig.Config, tags []string, dry bool) error {
	return nil
}
func (s stubLister) List(ctx context.Context) []plugin.CommandInfo { return s.commands }

func TestList(t *testing.T) {
	want := []plugin.CommandInfo{{FullPath: "welcome"}}

	got, err := plugin.List(t.Context(), stubLister{commands: want})
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestList_NotALister(t *testing.T) {
	_, err := plugin.List(t.Context(), struct{}{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not support the agent command catalog")
}

// stubSkillLister opts into both new interfaces.
type stubSkillLister struct {
	stubLister

	skill []plugin.SkillCommand
	meta  plugin.SkillMetadata
}

func (s stubSkillLister) ListSkill(ctx context.Context) []plugin.SkillCommand { return s.skill }
func (s stubSkillLister) SkillMetadata(ctx context.Context) plugin.SkillMetadata {
	return s.meta
}

func TestListSkill(t *testing.T) {
	want := []plugin.SkillCommand{
		{CommandInfo: plugin.CommandInfo{FullPath: "welcome"}, Skill: "#### Notes"},
	}

	got, err := plugin.ListSkill(t.Context(), stubSkillLister{skill: want})
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

// TestListSkill_ListerFallback covers the compatibility path: a plugin that
// only implements Lister still yields a catalog, just without prose.
func TestListSkill_ListerFallback(t *testing.T) {
	got, err := plugin.ListSkill(t.Context(), stubLister{
		commands: []plugin.CommandInfo{{FullPath: "welcome", Description: "greet"}},
	})
	require.NoError(t, err)

	require.Len(t, got, 1)
	assert.Equal(t, "welcome", got[0].FullPath)
	assert.Empty(t, got[0].Skill)
}

func TestListSkill_NotALister(t *testing.T) {
	_, err := plugin.ListSkill(t.Context(), struct{}{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not support the agent command catalog")
}

func TestSkillMetadataOf(t *testing.T) {
	want := plugin.SkillMetadata{Name: "acme"}

	assert.Equal(t, want, plugin.SkillMetadataOf(t.Context(), stubSkillLister{meta: want}))

	// A plugin not implementing SkillMetadataer falls back to the defaults.
	assert.Equal(t, plugin.SkillMetadata{}, plugin.SkillMetadataOf(t.Context(), stubLister{}))
}

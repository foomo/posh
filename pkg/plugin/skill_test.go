package plugin_test

import (
	"context"
	"fmt"
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

// kubectl is a command with nested structure, used to check what the renderers
// do and do not carry over from it.
func kubectl() plugin.SkillCommand {
	return plugin.SkillCommand{
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
	}
}

// TestRenderRootSkill covers the root skill: it indexes every command, including
// the ones that get no skill of their own, and points at the skills that exist.
func TestRenderRootSkill(t *testing.T) {
	out := plugin.RenderRootSkill(plugin.SkillMetadata{}, []plugin.SkillCommand{
		kubectl(),
		{CommandInfo: plugin.CommandInfo{FullPath: "exit", Description: "exit shell"}},
	})

	assert.Contains(t, out, "name: posh")

	// Every command is indexed with its description, whether or not it has a
	// skill of its own - the index is the only place a bare command appears.
	assert.Contains(t, out, "- `kubectl` - manage kubernetes clusters")
	assert.Contains(t, out, "- `exit` - exit shell")

	// A command with a skill of its own names it, so an agent can find it.
	assert.Contains(t, out, "(skill: `posh-kubectl`)")
	assert.NotContains(t, out, "- `exit` - exit shell (skill:")
}

// TestRenderRootSkill_NoDetail is the point of the split: the root carries no
// per-command detail at all. The flag and argument tree was over a quarter of
// the previous single file, and every leaf repeated the flags it inherited.
func TestRenderRootSkill_NoDetail(t *testing.T) {
	out := plugin.RenderRootSkill(plugin.SkillMetadata{}, []plugin.SkillCommand{kubectl()})

	assert.NotContains(t, out, "--dry-run", "flags must not reach the root skill")
	assert.NotContains(t, out, "(bool)")
	assert.NotContains(t, out, "manifest directory", "arguments must not reach the root skill")
	assert.NotContains(t, out, "kubectl <cluster> apply", "subcommand paths belong to the command's own skill")
}

// TestRenderRootSkill_Budget pins the size the split exists to buy. The root is
// loaded unconditionally, so it has to stay worth loading; the previous single
// file reached ~60k tokens in a real project.
func TestRenderRootSkill_Budget(t *testing.T) {
	commands := make([]plugin.SkillCommand, 0, 40)
	for i := range 40 {
		commands = append(commands, plugin.SkillCommand{
			CommandInfo: plugin.CommandInfo{
				FullPath:    fmt.Sprintf("command-%d", i),
				Description: "does something to the project",
				Flags:       []plugin.FlagInfo{{Name: "dry-run", Type: "bool"}},
			},
		})
	}

	// The fixed prose is the floor, and the budget the split was designed to:
	// everything an agent pays for unconditionally.
	fixed := len(strings.Fields(plugin.RenderRootSkill(plugin.SkillMetadata{}, nil)))
	assert.Less(t, fixed, 400, "the root skill's fixed prose must stay under 400 words")

	// Above that floor a command costs one index line, so growth is linear and
	// small - the single file it replaces grew by a command's whole flag tree.
	// At the 38 commands of the project that motivated this, ~700 words total.
	words := len(strings.Fields(plugin.RenderRootSkill(plugin.SkillMetadata{}, commands)))
	assert.Less(t, (words-fixed)/len(commands), 15, "a command must cost the root skill about one line")
}

// TestRenderRootSkill_Setup covers posh's own subcommands, which are documented
// from a fixed string rather than the catalog - so they render even for a project
// whose plugin registers no commands at all.
func TestRenderRootSkill_Setup(t *testing.T) {
	out := plugin.RenderRootSkill(plugin.SkillMetadata{}, nil)

	assert.Contains(t, out, "## Setup")
	assert.Contains(t, out, "`.posh.yaml`")
	assert.Contains(t, out, "`posh require`")
	assert.Contains(t, out, "`posh brew`")
	assert.Contains(t, out, "`posh prompt`")

	// The access control disclaimer is not optional.
	assert.Contains(t, out, "does not enforce access control")
}

// TestRenderRootSkill_Conventions covers the hoisted boilerplate: stated once
// here rather than repeated in every command's skill.
func TestRenderRootSkill_Conventions(t *testing.T) {
	out := plugin.RenderRootSkill(plugin.SkillMetadata{}, nil)

	assert.Contains(t, out, "## Conventions")
	assert.Contains(t, out, "`$id`")
	assert.Contains(t, out, "overridable via `.posh.yaml`")
}

// TestRenderCommandSkill covers a command's own skill: its runnable paths and
// its prose, but never its flags.
func TestRenderCommandSkill(t *testing.T) {
	out := plugin.RenderCommandSkill(kubectl())

	assert.Contains(t, out, "name: posh-kubectl")
	assert.Contains(t, out, "# `kubectl`")
	assert.Contains(t, out, "manage kubernetes clusters")

	// The runnable leaf paths are what a flat listing cannot show, so they stay.
	assert.Contains(t, out, "- `kubectl <cluster>` - target cluster")
	assert.Contains(t, out, "  - `kubectl <cluster> apply [Path]...` - apply manifests")

	// The detail they used to carry is replaced by a pointer to it.
	assert.NotContains(t, out, "--dry-run")
	assert.NotContains(t, out, "manifest directory")
	assert.Contains(t, out, "`posh agent catalog`")
	assert.Contains(t, out, "`posh help <command>`")
}

// TestRenderCommandSkill_Skill covers the command.Skiller contribution: appended
// verbatim, after the structure.
func TestRenderCommandSkill_Skill(t *testing.T) {
	out := plugin.RenderCommandSkill(plugin.SkillCommand{
		CommandInfo: plugin.CommandInfo{FullPath: "cache", Description: "manage caches"},
		// The trailing newline must not stack onto the one the renderer adds.
		Skill: "#### Configuration\n\n```yaml\ncache:\n  ttl: 5m\n```\n",
	})

	assert.Contains(t, out, "#### Configuration")
	assert.Contains(t, out, "  ttl: 5m")
	assert.NotContains(t, out, "```\n\n\n", "a trailing newline must not stack blank lines")
}

// TestRenderCommandSkill_RegisteredName covers the bug the Skiller signature
// change exists for: a provider registered under a name other than its default
// must document the name an agent actually types.
func TestRenderCommandSkill_RegisteredName(t *testing.T) {
	// What a provider whose default name is "squadron" renders once it is
	// registered as "admiral" and writes its prose against the name it is given.
	out := plugin.RenderCommandSkill(plugin.SkillCommand{
		CommandInfo: plugin.CommandInfo{FullPath: "admiral", Description: "manage squadrons"},
		Skill:       "#### Notes\n\nRun `posh x admiral up` to deploy.",
	})

	assert.Contains(t, out, "name: posh-admiral")
	assert.Contains(t, out, "posh x admiral up")
	assert.Contains(t, out, "posh execute admiral")
	assert.NotContains(t, out, "squadron up", "the default name must not leak into the invocation")
}

// TestRenderCommandSkill_Frontmatter covers the per-command overrides, including
// the YAML-hostile description that hand-built frontmatter would corrupt.
func TestRenderCommandSkill_Frontmatter(t *testing.T) {
	out := plugin.RenderCommandSkill(plugin.SkillCommand{
		CommandInfo: plugin.CommandInfo{FullPath: "squadron", Description: "manage squadrons"},
		Metadata: plugin.SkillMetadata{
			Name:         "acme-squadron",
			Description:  "Use when deploying: rollout, seed the DB #fast",
			AllowedTools: []string{"Bash(posh execute:*)"},
		},
	})

	front, _, ok := strings.Cut(strings.TrimPrefix(out, "---\n"), "---\n")
	require.True(t, ok, "the frontmatter must be delimited")

	var actual plugin.SkillMetadata
	require.NoError(t, yaml.Unmarshal([]byte(front), &actual), "the frontmatter must be valid YAML")

	assert.Equal(t, "acme-squadron", actual.Name)
	assert.Equal(t, "Use when deploying: rollout, seed the DB #fast", actual.Description)
	assert.Equal(t, []string{"Bash(posh execute:*)"}, actual.AllowedTools)

	// The hyphen form is what Claude Code reads; allowed_tools is not accepted.
	assert.Contains(t, front, "allowed-tools:")
}

// TestRenderCommandSkill_DescriptionFallback covers the derived description: it
// is a fallback, and a poor one, which is why install reports it.
func TestRenderCommandSkill_DescriptionFallback(t *testing.T) {
	out := plugin.RenderCommandSkill(plugin.SkillCommand{
		CommandInfo: plugin.CommandInfo{FullPath: "cache", Description: "manage caches"},
		Skill:       "#### Notes\n\nIn-memory.",
	})

	assert.Contains(t, out, "Use when running `cache` commands in this project.")
	assert.Contains(t, out, "Manage caches.", "the one-line description is sentence-cased into the fallback")
}

// TestRenderRootSkill_FrontmatterDefaults covers the zero value, with no empty
// allowed-tools key. The default description names triggering conditions: it is
// all a runtime sees when deciding whether to load the skill, so one that only
// says what posh is would never fire.
func TestRenderRootSkill_FrontmatterDefaults(t *testing.T) {
	out := plugin.RenderRootSkill(plugin.SkillMetadata{}, nil)

	assert.Contains(t, out, "name: posh")
	assert.Contains(t, out, "Use when")
	assert.NotContains(t, out, "allowed-tools")
}

// TestDescribes covers which commands earn a skill of their own. Arguments and
// flags do not count: a skill never prints them, so such a file would say
// nothing the root index does not already say.
func TestDescribes(t *testing.T) {
	for name, tt := range map[string]struct {
		command plugin.SkillCommand
		want    bool
	}{
		"prose": {
			command: plugin.SkillCommand{Skill: "#### Notes\n\nWorth saying."},
			want:    true,
		},
		"subcommands": {
			command: plugin.SkillCommand{CommandInfo: plugin.CommandInfo{
				Subcommands: []plugin.CommandInfo{{FullPath: "cache clear"}},
			}},
			want: true,
		},
		// help's shape: a hand-built CommandInfo carrying one optional argument.
		"arguments only": {
			command: plugin.SkillCommand{CommandInfo: plugin.CommandInfo{
				FullPath:  "help",
				Arguments: []plugin.ArgInfo{{Name: "command", Optional: true}},
			}},
			want: false,
		},
		"flags only": {
			command: plugin.SkillCommand{CommandInfo: plugin.CommandInfo{
				Flags: []plugin.FlagInfo{{Name: "dry-run", Type: "bool"}},
			}},
			want: false,
		},
		"bare leaf": {
			command: plugin.SkillCommand{CommandInfo: plugin.CommandInfo{FullPath: "exit"}},
			want:    false,
		},
		"whitespace prose": {
			command: plugin.SkillCommand{Skill: "  \n\t\n"},
			want:    false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.command.Describes())
		})
	}
}

// TestFallbackSkillDescriptions covers the report behind install's warning: only
// commands that get a skill and did not describe it are named.
func TestFallbackSkillDescriptions(t *testing.T) {
	got := plugin.FallbackSkillDescriptions([]plugin.SkillCommand{
		// Gets a skill, no description: reported.
		{CommandInfo: plugin.CommandInfo{FullPath: "cache"}, Skill: "#### Notes"},
		// Gets a skill and described it: fine.
		{
			CommandInfo: plugin.CommandInfo{FullPath: "squadron"},
			Skill:       "#### Notes",
			Metadata:    plugin.SkillMetadata{Description: "Use when deploying."},
		},
		// Gets no skill of its own, so it has nothing to describe.
		{CommandInfo: plugin.CommandInfo{FullPath: "exit"}},
	})

	assert.Equal(t, []string{"cache"}, got)
}

// TestWriteSkill covers the layout: a root skill plus one per command that has
// something to say, and none for the commands that do not.
func TestWriteSkill(t *testing.T) {
	dir := t.TempDir()

	written, err := plugin.WriteSkill(dir, plugin.SkillMetadata{}, []plugin.SkillCommand{
		{
			CommandInfo: plugin.CommandInfo{FullPath: "welcome", Description: "print a welcome message"},
			Skill:       "#### Notes\n\nNo side effects.",
		},
		{CommandInfo: plugin.CommandInfo{FullPath: "exit", Description: "exit shell"}},
	})
	require.NoError(t, err)

	assert.Equal(t, []string{
		filepath.Join(dir, "posh", "SKILL.md"),
		filepath.Join(dir, "posh-welcome", "SKILL.md"),
	}, written, "a command with nothing to say gets no skill of its own")

	b, err := os.ReadFile(filepath.Join(dir, "posh-welcome", "SKILL.md"))
	require.NoError(t, err)

	assert.Contains(t, string(b), "name: posh-welcome")
	assert.Contains(t, string(b), "No side effects.")

	// The root indexes both, including the one without its own skill.
	b, err = os.ReadFile(filepath.Join(dir, "posh", "SKILL.md"))
	require.NoError(t, err)

	assert.Contains(t, string(b), "- `welcome`")
	assert.Contains(t, string(b), "- `exit`")
}

func TestWriteSkill_DefaultPath(t *testing.T) {
	// DefaultSkillsPath is relative, so write it inside a temp working dir.
	dir := t.TempDir()
	t.Chdir(dir)

	_, err := plugin.WriteSkill("", plugin.SkillMetadata{}, nil)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(dir, plugin.DefaultSkillsPath, "posh", "SKILL.md"))
	assert.NoError(t, err, "an empty path must fall back to DefaultSkillsPath")
}

// TestWriteSkill_PrunesStale is why WriteSkill removes before it writes: nothing
// ever revisits the skill of a command that was renamed or dropped, so without
// this it would survive indefinitely and keep telling an agent to run something
// that no longer exists.
func TestWriteSkill_PrunesStale(t *testing.T) {
	dir := t.TempDir()

	_, err := plugin.WriteSkill(dir, plugin.SkillMetadata{}, []plugin.SkillCommand{
		{CommandInfo: plugin.CommandInfo{FullPath: "squadron"}, Skill: "#### Notes"},
	})
	require.NoError(t, err)

	require.FileExists(t, filepath.Join(dir, "posh-squadron", "SKILL.md"))

	// The same provider, registered under a different name this time.
	_, err = plugin.WriteSkill(dir, plugin.SkillMetadata{}, []plugin.SkillCommand{
		{CommandInfo: plugin.CommandInfo{FullPath: "admiral"}, Skill: "#### Notes"},
	})
	require.NoError(t, err)

	assert.NoDirExists(t, filepath.Join(dir, "posh-squadron"), "the renamed command's skill must be gone")
	assert.FileExists(t, filepath.Join(dir, "posh-admiral", "SKILL.md"))
}

// TestRemoveSkill covers uninstall, and that it stays inside what posh
// generated: the skills directory is shared with hand-written skills.
func TestRemoveSkill(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	root := filepath.Join(dir, plugin.DefaultSkillsPath)

	_, err := plugin.WriteSkill("", plugin.SkillMetadata{}, []plugin.SkillCommand{
		{CommandInfo: plugin.CommandInfo{FullPath: "cache"}, Skill: "#### Notes"},
	})
	require.NoError(t, err)

	// A skill posh did not generate, sitting alongside the ones it did.
	mine := filepath.Join(root, "my-own-skill")
	require.NoError(t, os.MkdirAll(mine, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(mine, "SKILL.md"), []byte("mine"), 0o644))

	require.NoError(t, plugin.RemoveSkill(""))

	assert.NoDirExists(t, filepath.Join(root, "posh"))
	assert.NoDirExists(t, filepath.Join(root, "posh-cache"))
	assert.FileExists(t, filepath.Join(mine, "SKILL.md"), "a hand-written skill must survive uninstall")

	assert.NoError(t, plugin.RemoveSkill(""), "removing missing skills must be a no-op")
	assert.NoError(t, plugin.RemoveSkill(filepath.Join(dir, "nonexistent")),
		"a missing skills directory must be a no-op")
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

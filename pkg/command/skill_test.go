package command_test

import (
	"context"
	"strings"
	"testing"

	"github.com/foomo/posh/pkg/cache"
	"github.com/foomo/posh/pkg/command"
	"github.com/foomo/posh/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// skiller is a command contributing extra SKILL.md markdown. It writes its prose
// against the name it is given, as a command registered under a non-default name
// must.
type skiller struct{ leaf }

func (skiller) Skill(ctx context.Context, name string) string {
	return "#### Notes\n\nRun `" + name + "` twice."
}

// metadataer is a command supplying its own skill frontmatter.
type metadataer struct{ leaf }

func (metadataer) SkillMetadata(ctx context.Context, name string) command.SkillMetadata {
	return command.SkillMetadata{Description: "Use when " + name + " is needed."}
}

// TestSkill_NotASkiller covers the fallback: a command that does not opt in
// contributes nothing.
func TestSkill_NotASkiller(t *testing.T) {
	assert.Empty(t, command.Skill(t.Context(), leaf{}, "leaf"))
}

// TestSkill_Skiller covers the dispatch, and that the registered name reaches
// the contribution: the same command can be registered under another name, and
// prose naming a command that does not exist in the project is worse than none.
func TestSkill_Skiller(t *testing.T) {
	assert.Equal(t, "#### Notes\n\nRun `admiral` twice.",
		command.Skill(t.Context(), skiller{}, "admiral"))
}

// TestSkillMetadataOf covers the per-command frontmatter, and the zero value a
// command that does not opt in falls back to.
func TestSkillMetadataOf(t *testing.T) {
	assert.Equal(t,
		command.SkillMetadata{Description: "Use when admiral is needed."},
		command.SkillMetadataOf(t.Context(), metadataer{}, "admiral"))

	assert.Equal(t, command.SkillMetadata{}, command.SkillMetadataOf(t.Context(), leaf{}, "leaf"))
}

// TestSkill_NoSelfDefeatingBuiltIns covers the two deliberate omissions. Both
// commands' prose would only steer an agent away from them - exit does nothing
// under `posh execute`, and the catalog beats help - and a skill has to be
// loaded before an agent can read that, so contributing none is cheaper.
func TestSkill_NoSelfDefeatingBuiltIns(t *testing.T) {
	l := log.NewTest(t)

	assert.Empty(t, command.Skill(t.Context(), command.NewExit(l), "exit"))
	assert.Empty(t, command.Skill(t.Context(), command.NewHelp(l, command.Commands{}), "help"))
}

// TestSkill_BuiltIns asserts the remaining built-ins contribute prose, and that
// each contribution is renderable: a heading above level 4 would collide with
// the command heading RenderCommandSkill puts above it.
func TestSkill_BuiltIns(t *testing.T) {
	l := log.NewTest(t)

	for name, cmd := range map[string]any{
		"cache":   command.NewCache(l, cache.NewMemoryCache()),
		"check":   command.NewCheck(l),
		"env":     command.NewEnv(l),
		"history": command.NewHistory(l, nil),
	} {
		t.Run(name, func(t *testing.T) {
			actual := command.Skill(t.Context(), cmd, name)

			require.NotEmpty(t, actual, "every built-in should contribute skill prose")
			assert.True(t, strings.HasPrefix(actual, "#### "),
				"contributions must start at heading level 4, below the command's own")

			for line := range strings.SplitSeq(actual, "\n") {
				assert.False(t, strings.HasPrefix(line, "### ") || strings.HasPrefix(line, "## "),
					"a heading above level 4 would collide with the command heading: %q", line)
			}
		})
	}
}

// TestDescribe_Help covers the hand-built CommandInfo: help is not tree based,
// so without it the optional [command] argument would be invisible.
func TestDescribe_Help(t *testing.T) {
	actual := command.NewHelp(log.NewTest(t), command.Commands{}).Describe(t.Context())

	assert.Equal(t, "help", actual.FullPath)

	require.Len(t, actual.Arguments, 1)
	assert.Equal(t, "command", actual.Arguments[0].Name)
	assert.True(t, actual.Arguments[0].Optional)

	// Usage is rebuilt from the arguments, so an optional one is angled.
	assert.Equal(t, "<command>", actual.Usage())
}

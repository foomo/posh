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

// skiller is a command contributing extra SKILL.md markdown.
type skiller struct{ leaf }

func (skiller) Skill(ctx context.Context) string { return "#### Notes\n\nRun it twice." }

// TestSkill_NotASkiller covers the fallback: a command that does not opt in
// contributes nothing.
func TestSkill_NotASkiller(t *testing.T) {
	assert.Empty(t, command.Skill(t.Context(), leaf{}))
}

// TestSkill_Skiller covers the dispatch.
func TestSkill_Skiller(t *testing.T) {
	assert.Equal(t, "#### Notes\n\nRun it twice.", command.Skill(t.Context(), skiller{}))
}

// TestSkill_BuiltIns asserts every built-in command contributes prose, and that
// each contribution is renderable: a "###" heading would collide with the
// command heading RenderSkill puts above it.
func TestSkill_BuiltIns(t *testing.T) {
	l := log.NewTest(t)

	for name, cmd := range map[string]any{
		"cache":   command.NewCache(l, cache.NewMemoryCache()),
		"check":   command.NewCheck(l),
		"env":     command.NewEnv(l),
		"exit":    command.NewExit(l),
		"help":    command.NewHelp(l, command.Commands{}),
		"history": command.NewHistory(l, nil),
	} {
		t.Run(name, func(t *testing.T) {
			actual := command.Skill(t.Context(), cmd)

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

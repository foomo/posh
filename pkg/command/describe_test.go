package command_test

import (
	"testing"

	"github.com/foomo/posh/pkg/cache"
	"github.com/foomo/posh/pkg/command"
	"github.com/foomo/posh/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// leaf is a command that does not describe itself.
type leaf struct{}

func (leaf) Name() string        { return "exit" }
func (leaf) Description() string { return "exit shell" }

// TestDescribe_Leaf covers the fallback: a command that is not a Describer is
// still listed, just without subcommand detail.
func TestDescribe_Leaf(t *testing.T) {
	actual := command.Describe(t.Context(), "exit", "exit shell", leaf{})

	assert.Equal(t, "exit", actual.FullPath)
	assert.Equal(t, "exit shell", actual.Description)
	assert.Empty(t, actual.Subcommands)
	assert.Empty(t, actual.Flags)
	assert.False(t, actual.Dynamic)
}

// TestDescribe_Describer covers the dispatch: a Describer describes itself, so
// the name and description arguments are ignored in favour of the tree's own.
func TestDescribe_Describer(t *testing.T) {
	l := log.NewTest(t)

	actual := command.Describe(t.Context(), "ignored", "ignored", command.NewCache(l, cache.NewMemoryCache()))

	assert.Equal(t, "cache", actual.FullPath)
	assert.Equal(t, "Manage the internal cache", actual.Description)

	require.NotEmpty(t, actual.Subcommands)
	assert.Equal(t, "cache clear", actual.Subcommands[0].FullPath)
}

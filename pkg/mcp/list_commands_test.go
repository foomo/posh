package mcp_test

import (
	"context"
	"testing"

	ownbrewconfig "github.com/foomo/ownbrew/pkg/config"
	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/mcp"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubPlugin implements plugin.Plugin and plugin.Lister, matching the stub
// pattern already used in pkg/plugin/skill_test.go.
type stubPlugin struct {
	commands []plugin.CommandInfo
}

func (s stubPlugin) Prompt(ctx context.Context, cfg config.Prompt) error   { return nil }
func (s stubPlugin) Execute(ctx context.Context, args []string) error      { return nil }
func (s stubPlugin) Require(ctx context.Context, cfg config.Require) error { return nil }
func (s stubPlugin) Brew(ctx context.Context, cfg ownbrewconfig.Config, tags []string, dry bool) error {
	return nil
}
func (s stubPlugin) List(ctx context.Context) []plugin.CommandInfo { return s.commands }

func TestListCommands(t *testing.T) {
	plg := stubPlugin{commands: []plugin.CommandInfo{{FullPath: "welcome", Description: "greet"}}}

	got, err := mcp.ListCommands(t.Context(), plg)
	require.NoError(t, err)

	assert.Equal(t, plugin.NewCatalog(plg.commands), got)
}

func TestListCommands_NotALister(t *testing.T) {
	_, err := mcp.ListCommands(t.Context(), struct{}{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not support the agent command catalog")
}

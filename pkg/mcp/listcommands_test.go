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

	got, err := mcp.ListCommands(t.Context(), plg, "", 0)
	require.NoError(t, err)

	assert.Equal(t, plugin.NewCatalog(plg.commands), got)
}

func TestListCommands_NotALister(t *testing.T) {
	_, err := mcp.ListCommands(t.Context(), struct{}{}, "", 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not support the agent command catalog")
}

func squadronTree() plugin.CommandInfo {
	return plugin.CommandInfo{
		FullPath: "squadron", Description: "manage squadrons",
		Subcommands: []plugin.CommandInfo{
			{FullPath: "squadron <cluster>", Dynamic: true, Subcommands: []plugin.CommandInfo{
				{FullPath: "squadron <cluster> <fleet>", Dynamic: true, Subcommands: []plugin.CommandInfo{
					{FullPath: "squadron <cluster> <fleet> up", Description: "installs a chart"},
					{FullPath: "squadron <cluster> <fleet> down", Description: "uninstalls a chart"},
				}},
			}},
		},
	}
}

func TestListCommands_Path(t *testing.T) {
	plg := stubPlugin{commands: []plugin.CommandInfo{
		{FullPath: "welcome", Description: "greet"},
		squadronTree(),
	}}

	got, err := mcp.ListCommands(t.Context(), plg, "squadron", 0)
	require.NoError(t, err)

	assert.Equal(t, []plugin.CommandInfo{squadronTree()}, got.Commands)
}

func TestListCommands_PathNotFound(t *testing.T) {
	plg := stubPlugin{commands: []plugin.CommandInfo{{FullPath: "welcome"}}}

	got, err := mcp.ListCommands(t.Context(), plg, "nope", 0)
	require.NoError(t, err)

	assert.Empty(t, got.Commands)
}

func TestListCommands_Depth(t *testing.T) {
	plg := stubPlugin{commands: []plugin.CommandInfo{squadronTree()}}

	got, err := mcp.ListCommands(t.Context(), plg, "", 1)
	require.NoError(t, err)

	require.Len(t, got.Commands, 1)
	assert.Equal(t, "squadron", got.Commands[0].FullPath)
	require.Len(t, got.Commands[0].Subcommands, 1)
	assert.Equal(t, "squadron <cluster>", got.Commands[0].Subcommands[0].FullPath)
	assert.Empty(t, got.Commands[0].Subcommands[0].Subcommands, "depth 1 must stop one level below the top")
}

func TestListCommands_PathAndDepth(t *testing.T) {
	plg := stubPlugin{commands: []plugin.CommandInfo{
		{FullPath: "welcome"},
		squadronTree(),
	}}

	got, err := mcp.ListCommands(t.Context(), plg, "squadron", 1)
	require.NoError(t, err)

	require.Len(t, got.Commands, 1)
	assert.Equal(t, "squadron", got.Commands[0].FullPath)
	require.Len(t, got.Commands[0].Subcommands, 1)
	assert.Empty(t, got.Commands[0].Subcommands[0].Subcommands)
}

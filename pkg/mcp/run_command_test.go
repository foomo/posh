package mcp_test

import (
	"context"
	"testing"

	ownbrewconfig "github.com/foomo/ownbrew/pkg/config"
	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/mcp"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// executingPlugin logs through whatever logger it is constructed with and
// optionally fails, matching how a real downstream plugin's Execute behaves.
type executingPlugin struct {
	l       log.Logger
	failMsg string
}

func (p executingPlugin) Prompt(ctx context.Context, cfg config.Prompt) error   { return nil }
func (p executingPlugin) Require(ctx context.Context, cfg config.Require) error { return nil }
func (p executingPlugin) Brew(ctx context.Context, cfg ownbrewconfig.Config, tags []string, dry bool) error {
	return nil
}
func (p executingPlugin) Execute(ctx context.Context, args []string) error {
	p.l.Info("running", args)

	if p.failMsg != "" {
		return errors.New(p.failMsg)
	}

	return nil
}

func TestRunCommand_Success(t *testing.T) {
	provider := func(l log.Logger) (plugin.Plugin, error) {
		return executingPlugin{l: l}, nil
	}

	output, err := mcp.RunCommand(t.Context(), provider, []string{"welcome"})
	require.NoError(t, err)

	assert.Contains(t, output, `"message":"running [welcome]"`)
}

func TestRunCommand_Failure(t *testing.T) {
	provider := func(l log.Logger) (plugin.Plugin, error) {
		return executingPlugin{l: l, failMsg: "boom"}, nil
	}

	output, err := mcp.RunCommand(t.Context(), provider, []string{"broken"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
	assert.Contains(t, output, `"message":"running [broken]"`)
}

func TestRunCommand_MissingArgs(t *testing.T) {
	provider := func(l log.Logger) (plugin.Plugin, error) {
		return executingPlugin{l: l}, nil
	}

	_, err := mcp.RunCommand(t.Context(), provider, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing [cmd] argument")
}

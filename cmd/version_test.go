package cmd_test

import (
	"encoding/json"
	"testing"

	"github.com/foomo/posh/cmd"
	intversion "github.com/foomo/posh/internal/version"
	"github.com/foomo/posh/pkg/agent"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func versionCmd(t *testing.T) *cobra.Command {
	t.Helper()

	root := &cobra.Command{Use: "posh"}
	cmd.NewVersion(root)

	for _, c := range root.Commands() {
		if c.Name() == "version" {
			return c
		}
	}

	t.Fatal("version command not registered")

	return nil
}

// TestVersion_AgentMode covers the regression this change is about: the version
// used to be delivered as a pkg/log.AgentJSON entry, so the result carried the
// "type" field that is supposed to mark an interleaved log line rather than a
// result value.
func TestVersion_AgentMode(t *testing.T) {
	defer agent.SetDetected(true)()

	c := versionCmd(t)

	var err error

	out := captureStdout(t, func() { err = c.RunE(c, nil) })
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &got), "agent output must be one JSON value")

	assert.Equal(t, intversion.Version, got["version"])
	assert.NotContains(t, got, "type", "a result value must not carry the log discriminator")
	assert.NotContains(t, got, "message")

	// Commit and build time are debug-level only, and the default level is info.
	assert.NotContains(t, got, "commit")
	assert.NotContains(t, got, "build_time")
}

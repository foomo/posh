package cmd_test

import (
	"encoding/json"
	"testing"

	"github.com/foomo/posh/cmd"
	"github.com/foomo/posh/pkg/agent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// configCmd registers the config command on a throwaway root and returns it.
//
// The command is driven through RunE directly rather than Execute: PreRunE
// loads .posh.yaml from the working directory, which is the config loader's
// behavior to test, not this command's output.
func configCmd(t *testing.T) *cobra.Command {
	t.Helper()

	root := &cobra.Command{Use: "posh"}
	cmd.NewConfig(root)

	for _, c := range root.Commands() {
		if c.Name() == "config" {
			return c
		}
	}

	t.Fatal("config command not registered")

	return nil
}

func TestConfig_AgentMode(t *testing.T) {
	defer agent.SetDetected(true)()

	viper.Set("some-key", "some-value")
	defer viper.Set("some-key", nil)

	c := configCmd(t)

	var err error

	out := captureStdout(t, func() { err = c.RunE(c, nil) })
	require.NoError(t, err)

	assert.NotContains(t, out, "\033", "agent output must carry no ANSI escapes")

	var got map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &got), "agent output must be one JSON value")
	assert.Equal(t, "some-value", got["some-key"])
}

// TestConfig_NoColor covers the --no-color flag reaching this command: it
// highlights via chroma, which the pterm global the flag otherwise toggles does
// not affect.
func TestConfig_NoColor(t *testing.T) {
	defer agent.SetDetected(false)()

	viper.Set("no-color", true)
	defer viper.Set("no-color", false)

	viper.Set("some-key", "some-value")
	defer viper.Set("some-key", nil)

	c := configCmd(t)

	var err error

	out := captureStdout(t, func() { err = c.RunE(c, nil) })
	require.NoError(t, err)

	assert.NotContains(t, out, "\033")
	assert.Contains(t, out, "some-key: some-value")
}

// TestConfig_Color is the counterpart: without the flag the human path still
// highlights, so the no-color assertion above is testing something.
func TestConfig_Color(t *testing.T) {
	defer agent.SetDetected(false)()

	viper.Set("some-key", "some-value")
	defer viper.Set("some-key", nil)

	c := configCmd(t)

	var err error

	out := captureStdout(t, func() { err = c.RunE(c, nil) })
	require.NoError(t, err)

	assert.Contains(t, out, "\033", "human output is syntax highlighted")
}

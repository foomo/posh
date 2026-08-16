package command_test

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/foomo/posh/pkg/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnv_ListAgent_KeysOnly(t *testing.T) {
	t.Setenv("POSH_TEST_ENV", "value-one")

	out := captureAgent(t, func() {
		l := newTestLogger()
		require.NoError(t, command.NewEnv(l).Execute(t.Context(), parse(t, l, "env list")))
	})

	var got command.EnvList
	require.NoError(t, json.Unmarshal([]byte(out), &got))

	var found bool

	for _, v := range got.Values {
		// Without --with-values the value is omitted for every entry.
		assert.Empty(t, v.Value)

		if v.Name == "POSH_TEST_ENV" {
			found = true
		}
	}

	assert.True(t, found, "expected POSH_TEST_ENV in the listing")
}

func TestEnv_ListAgent_WithValues(t *testing.T) {
	// A value longer than any terminal width: the PTerm path wraps values to
	// the terminal width, which would corrupt this one with newlines.
	long := strings.Repeat("abcdefghij", 60)
	t.Setenv("POSH_TEST_LONG", long)

	out := captureAgent(t, func() {
		l := newTestLogger()
		require.NoError(t, command.NewEnv(l).Execute(t.Context(), parse(t, l, "env list --with-values")))
	})

	var got command.EnvList
	require.NoError(t, json.Unmarshal([]byte(out), &got))

	var found bool

	for _, v := range got.Values {
		if v.Name == "POSH_TEST_LONG" {
			found = true

			assert.Equal(t, long, v.Value, "value must not be wrapped or truncated")
		}
	}

	assert.True(t, found, "expected POSH_TEST_LONG in the listing")
}

func TestEnv_GetAgent(t *testing.T) {
	t.Setenv("POSH_TEST_ENV", "value-one")

	out := captureAgent(t, func() {
		l := newTestLogger()
		require.NoError(t, command.NewEnv(l).Execute(t.Context(), parse(t, l, "env get POSH_TEST_ENV")))
	})

	var got command.EnvGet
	require.NoError(t, json.Unmarshal([]byte(out), &got))

	assert.Equal(t, "POSH_TEST_ENV", got.Name)
	assert.Equal(t, "value-one", got.Value)
}

func TestEnv_SetUnset(t *testing.T) {
	ctx := context.Background()
	l := newTestLogger()
	c := command.NewEnv(l)

	require.NoError(t, c.Execute(ctx, parse(t, l, "env set POSH_TEST_SET has-value")))
	assert.Equal(t, "has-value", os.Getenv("POSH_TEST_SET"))

	require.NoError(t, c.Execute(ctx, parse(t, l, "env unset POSH_TEST_SET")))

	_, ok := os.LookupEnv("POSH_TEST_SET")
	assert.False(t, ok)
}

// TestEnv_ListHuman verifies the non-agent path renders via PTerm rather than
// emitting JSON. PTerm writes through its own cached writer instead of the
// os.Stdout we swap, so its output does not reach the pipe - the assertion is
// that no JSON envelope was produced.
func TestEnv_ListHuman(t *testing.T) {
	t.Setenv("POSH_TEST_ENV", "value-one")

	l := newTestLogger()

	out := captureHuman(t, func() {
		require.NoError(t, command.NewEnv(l).Execute(t.Context(), parse(t, l, "env list")))
	})

	// Results are emitted bare, so there is no envelope to key on: assert the
	// output is not the JSON payload at all.
	assert.NotContains(t, out, `"values"`)
}

func TestEnv_GetHuman(t *testing.T) {
	t.Setenv("POSH_TEST_ENV", "value-one")

	out := captureHuman(t, func() {
		l := newTestLogger()
		require.NoError(t, command.NewEnv(l).Execute(t.Context(), parse(t, l, "env get POSH_TEST_ENV")))
	})

	// The logger writes to os.Stdout directly, so this one is capturable: the
	// value arrives wrapped in a posh.log envelope rather than as a bare
	// EnvGet payload.
	assert.Contains(t, out, "value-one")
	assert.Contains(t, out, "posh.log")
	assert.NotContains(t, out, `"name"`)
}

func TestEnv_GetAgent_NoHTMLEscaping(t *testing.T) {
	t.Setenv("POSH_TEST_HTML", "a&b<c>d")

	out := captureAgent(t, func() {
		l := newTestLogger()
		require.NoError(t, command.NewEnv(l).Execute(t.Context(), parse(t, l, "env get POSH_TEST_HTML")))
	})

	assert.Contains(t, out, "a&b<c>d", "characters must not be escaped to \\u0026 etc")
}

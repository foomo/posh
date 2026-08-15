package command_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/foomo/posh/pkg/command"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt/check"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheck_RunAgent(t *testing.T) {
	// Two checkers, deliberately out of alphabetical order.
	zulu := func(ctx context.Context, l log.Logger) []check.Info {
		return []check.Info{check.NewFailureInfo("⚡︎", "zulu", "down")}
	}
	alpha := func(ctx context.Context, l log.Logger) []check.Info {
		return []check.Info{check.NewSuccessInfo("⚫︎", "alpha", "ok")}
	}

	out := captureAgent(t, func() {
		l := newTestLogger()
		require.NoError(t, command.NewCheck(l, zulu, alpha).Execute(t.Context(), parse(t, l, "check")))
	})

	var got check.Results
	require.NoError(t, json.Unmarshal([]byte(out), &got))

	require.Len(t, got.Checks, 2)

	// Sorted by name, matching the PTerm path.
	assert.Equal(t, "alpha", got.Checks[0].Name)
	assert.Equal(t, "ok", got.Checks[0].Note)
	assert.Equal(t, "success", got.Checks[0].Status)

	assert.Equal(t, "zulu", got.Checks[1].Name)
	assert.Equal(t, "failure", got.Checks[1].Status)

	// Icons are presentation and must not leak into the payload.
	assert.NotContains(t, out, "⚫︎")
}

func TestCheck_RunAgent_NoCheckers(t *testing.T) {
	out := captureAgent(t, func() {
		l := newTestLogger()
		require.NoError(t, command.NewCheck(l).Execute(t.Context(), parse(t, l, "check")))
	})

	var got check.Results
	require.NoError(t, json.Unmarshal([]byte(out), &got))

	assert.NotNil(t, got.Checks, "an empty result must encode as [] rather than null")
	assert.Empty(t, got.Checks)
}

// TestCheck_RunHuman verifies the non-agent path still goes through check.Check.
func TestCheck_RunHuman(t *testing.T) {
	var called bool

	l := newTestLogger()

	// NewCheck picks its check.Check at construction, so it must be built
	// inside captureHuman - while agent mode is off.
	out := captureHuman(t, func() {
		c := command.NewCheck(l, func(ctx context.Context, l log.Logger) []check.Info {
			called = true
			return []check.Info{check.NewSuccessInfo("⚫︎", "alpha", "ok")}
		})

		require.NoError(t, c.Execute(t.Context(), parse(t, l, "check")))
	})

	assert.True(t, called, "the checker still runs through check.Check")
	assert.NotContains(t, out, `"checks"`, "no JSON payload on the human path")
}

package agent_test

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/foomo/posh/pkg/agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender_AgentMode(t *testing.T) {
	restore := agent.SetDetected(true)
	defer restore()

	var humanCalled bool

	out := captureStdout(t, func() {
		require.NoError(t, agent.Render(
			func() any { return map[string]string{"key": "value"} },
			func() error { humanCalled = true; return nil },
		))
	})

	assert.False(t, humanCalled, "the human closure must not run in agent mode")
	assert.JSONEq(t, `{"key":"value"}`, out)
}

func TestRender_HumanMode(t *testing.T) {
	restore := agent.SetDetected(false)
	defer restore()

	var payloadCalled bool

	out := captureStdout(t, func() {
		require.NoError(t, agent.Render(
			func() any { payloadCalled = true; return nil },
			func() error { return nil },
		))
	})

	assert.False(t, payloadCalled, "the payload must not be built in human mode")
	assert.Empty(t, out)
}

func TestRender_HumanErrorPropagates(t *testing.T) {
	restore := agent.SetDetected(false)
	defer restore()

	want := errors.New("boom")
	got := agent.Render(func() any { return nil }, func() error { return want })

	assert.Equal(t, want, got)
}

// captureStdout swaps the process stdout for a pipe. The JSON encoder writes to
// the real descriptor, so a writer-only capture would miss it.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	orig := os.Stdout
	os.Stdout = w

	done := make(chan string, 1)

	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()

	defer func() { os.Stdout = orig }()

	fn()

	require.NoError(t, w.Close())

	return <-done
}

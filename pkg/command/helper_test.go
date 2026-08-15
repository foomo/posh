package command_test

import (
	"io"
	"os"
	"testing"

	"github.com/foomo/posh/pkg/agent"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/readline"
	"github.com/stretchr/testify/require"
)

func newTestLogger() log.Logger {
	return log.NewAgentJSON(log.AgentJSONWithLevel(log.LevelInfo))
}

// parse builds a Readline from a command line, wiring up the flag sets the way
// the prompt does so that flags such as --with-values are available.
func parse(t *testing.T, l log.Logger, input string) *readline.Readline {
	t.Helper()

	r, err := readline.New(l)
	require.NoError(t, err)
	require.NoError(t, r.Parse(input))

	return r
}

// captureStdout swaps the process stdout file descriptor for a pipe and returns
// everything written while fn runs.
//
// It captures the real descriptor rather than a logger writer because PTerm
// writes to os.Stdout directly: a writer-only capture would miss the
// human-formatted output entirely and silently pass.
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

	// Restore stdout and drain the pipe even if fn panics.
	defer func() {
		os.Stdout = orig
	}()

	fn()

	require.NoError(t, w.Close())

	return <-done
}

// captureAgent forces agent mode for the duration of fn and returns what was
// written to stdout.
func captureAgent(t *testing.T, fn func()) string {
	t.Helper()

	restore := agent.SetDetected(true)
	defer restore()

	return captureStdout(t, fn)
}

// captureHuman forces human mode for the duration of fn and returns what was
// written to stdout. Needed because the test binary itself usually runs inside
// an agent harness.
func captureHuman(t *testing.T, fn func()) string {
	t.Helper()

	restore := agent.SetDetected(false)
	defer restore()

	return captureStdout(t, fn)
}

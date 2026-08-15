package cmd_test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

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

package exec_test

import (
	"bytes"
	"testing"

	"github.com/foomo/posh/pkg/exec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommand_DefaultsToRealStdio(t *testing.T) {
	cmd := exec.NewCommand(t.Context(), "echo", "hi")

	assert.NotNil(t, cmd)
}

func TestNewCommand_UsesContextStdio(t *testing.T) {
	var stdout, stderr bytes.Buffer

	ctx := exec.WithStdio(t.Context(), &stdout, &stderr)

	err := exec.NewCommand(ctx, "echo", "hello-from-context").Run()
	require.NoError(t, err)

	assert.Contains(t, stdout.String(), "hello-from-context")
}

func TestStdioFrom_NotSet(t *testing.T) {
	_, _, ok := exec.StdioFrom(t.Context())
	assert.False(t, ok)
}

package prompt_test

import (
	"testing"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWithThemeUnconfigured pins the opt-in contract: an empty value is not an
// error and leaves the prompt exactly as it was before theming existed.
func TestWithThemeUnconfigured(t *testing.T) {
	t.Parallel()

	subject, err := prompt.New(log.NewFmt(), prompt.WithTheme(""))
	require.NoError(t, err)
	assert.Empty(t, prompt.ThemeOptions(subject), "an unconfigured prompt must add no color options")
}

func TestWithThemeConfigured(t *testing.T) {
	t.Parallel()

	subject, err := prompt.New(log.NewFmt(), prompt.WithTheme("latte"))
	require.NoError(t, err)
	assert.Len(t, prompt.ThemeOptions(subject), 13)
}

// TestWithThemeRejectsUnknown surfaces a typo in .posh.yaml at startup rather
// than rendering a different theme than the one asked for.
func TestWithThemeRejectsUnknown(t *testing.T) {
	t.Parallel()

	_, err := prompt.New(log.NewFmt(), prompt.WithTheme("nord"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nord")
}

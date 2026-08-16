// Tests in this package mutate gookit's detected color level, so they save and
// restore it in t.Cleanup and must not call t.Parallel — a leaked mutation
// surfaces as a confusing failure in an unrelated package.
package theme_test

import (
	"regexp"
	"testing"

	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/theme"
	"github.com/gookit/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var hexRe = regexp.MustCompile(`^#[0-9a-f]{6}$`)

// forceLevel pins gookit's color level for the duration of a test.
func forceLevel(t *testing.T, level color.Level) {
	t.Helper()

	prevLevel, prevEnable := color.TermColorLevel(), color.Enable
	color.Enable = true

	color.ForceSetColorLevel(level)

	t.Cleanup(func() {
		color.ForceSetColorLevel(prevLevel)

		color.Enable = prevEnable
	})
}

func TestPalettesAreComplete(t *testing.T) {
	forceLevel(t, color.Level256)

	for _, flavor := range theme.Flavors() {
		t.Run(flavor.String(), func(t *testing.T) {
			subject, err := theme.New(theme.WithFlavor(flavor))
			require.NoError(t, err)

			for _, role := range theme.Roles() {
				assert.Regexp(t, hexRe, subject.Hex(role), "role %s has no hex", role)
				assert.NotEqual(t, "unknown", role.String())
			}
		})
	}
}

// TestPromptColorsAreComplete covers the tier the prompt actually renders: a
// missing entry would silently fall back to the terminal default.
func TestPromptColorsAreComplete(t *testing.T) {
	forceLevel(t, color.Level256)

	for _, flavor := range theme.Flavors() {
		t.Run(flavor.String(), func(t *testing.T) {
			subject, err := theme.New(theme.WithFlavor(flavor))
			require.NoError(t, err)

			for _, role := range theme.Roles() {
				assert.NotEqual(t, goprompt.DefaultColor, subject.Prompt(role),
					"role %s falls back to the default color", role)
			}
		})
	}
}

// TestPromptColorsMatchCatppuccinANSI pins the 16-color mapping to Catppuccin's
// own published ANSI slots (the `ansiColors` block of catppuccin/palette).
//
// That mapping puts each palette's signature colors in the *normal* SGR slots
// and variant shades in the bright ones, so a role mapped to a bright slot
// renders a different shade than intended on a Catppuccin terminal — which is
// exactly the setup this tier exists to serve. An earlier version of the table
// used bright slots throughout for the dark flavors and was one shade off on
// six of the nine roles.
func TestPromptColorsMatchCatppuccinANSI(t *testing.T) {
	forceLevel(t, color.Level256)

	// The roles whose Catppuccin color IS the normal ANSI slot for that hue,
	// so the expectation is exact rather than nearest-neighbour.
	exact := map[theme.Role]goprompt.Color{
		theme.RoleSuccess: goprompt.DarkGreen, // green
		theme.RoleWarning: goprompt.Brown,     // yellow
		theme.RoleError:   goprompt.DarkRed,   // red
	}

	for _, flavor := range theme.Flavors() {
		t.Run(flavor.String(), func(t *testing.T) {
			subject, err := theme.New(theme.WithFlavor(flavor))
			require.NoError(t, err)

			for role, want := range exact {
				assert.Equal(t, want, subject.Prompt(role),
					"role %s must use the normal ANSI slot so a Catppuccin terminal renders the exact palette color", role)
			}
		})
	}
}

// TestFlavorsDiffer pins the user-visible promise: a light flavor must not
// render like a dark one.
func TestFlavorsDiffer(t *testing.T) {
	forceLevel(t, color.Level256)

	mocha, err := theme.New(theme.WithFlavor(theme.FlavorMocha))
	require.NoError(t, err)

	latte, err := theme.New(theme.WithFlavor(theme.FlavorLatte))
	require.NoError(t, err)

	assert.True(t, mocha.IsDark())
	assert.False(t, latte.IsDark())

	for _, role := range []theme.Role{theme.RoleText, theme.RolePrimary, theme.RoleError, theme.RoleMuted} {
		assert.NotEqual(t, mocha.Hex(role), latte.Hex(role), "role %s should differ between flavors", role)
	}

	// The text role in particular has to flip, or the input line is illegible
	// on one of the two backgrounds.
	assert.NotEqual(t, mocha.Prompt(theme.RoleText), latte.Prompt(theme.RoleText))
}

func TestNoColorYieldsDefaults(t *testing.T) {
	forceLevel(t, color.Level256)

	subject, err := theme.New(theme.WithFlavor(theme.FlavorMocha), theme.WithCapability(theme.CapabilityNone))
	require.NoError(t, err)

	for _, role := range theme.Roles() {
		assert.Equal(t, goprompt.DefaultColor, subject.Prompt(role))
	}
}

func TestParseFlavor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    theme.Flavor
		wantErr bool
	}{
		{name: "lowercase", input: "mocha", want: theme.FlavorMocha},
		{name: "mixed case", input: "Latte", want: theme.FlavorLatte},
		{name: "surrounding space", input: "  frappe  ", want: theme.FlavorFrappe},
		{name: "macchiato", input: "macchiato", want: theme.FlavorMacchiato},
		{name: "unknown", input: "solarized", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := theme.ParseFlavor(tt.input)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestResolveUnconfigured is the whole opt-in contract: no flavor means no
// theme, which callers use to keep their previous colors untouched.
func TestResolveUnconfigured(t *testing.T) {
	forceLevel(t, color.Level256)

	subject, err := theme.Resolve("")
	require.NoError(t, err)
	assert.Nil(t, subject)
}

func TestResolveConfigured(t *testing.T) {
	forceLevel(t, color.Level256)

	subject, err := theme.Resolve("latte")
	require.NoError(t, err)
	require.NotNil(t, subject)
	assert.Equal(t, theme.FlavorLatte, subject.Flavor())
}

// TestResolveRejectsUnknown covers a typo in something the user wrote: it must
// surface rather than silently render a different theme.
func TestResolveRejectsUnknown(t *testing.T) {
	forceLevel(t, color.Level256)

	_, err := theme.Resolve("nord")
	require.Error(t, err)
}

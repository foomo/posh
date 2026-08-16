// Package theme provides the Catppuccin palettes used to color the
// interactive prompt.
//
// It is opt-in: a project selects a flavor via the `prompt.theme` key in
// .posh.yaml, and without it posh keeps its previous, unthemed colors.
package theme

import (
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/pkg/errors"
)

type (
	// Theme is a resolved palette: a flavor plus the terminal capability its
	// colors have been projected onto.
	Theme struct {
		flavor     Flavor
		capability Capability
		hex        map[Role]string
	}
	Option func(*Theme) error
)

// ------------------------------------------------------------------------------------------------
// ~ Options
// ------------------------------------------------------------------------------------------------

func WithFlavor(v Flavor) Option {
	return func(o *Theme) error {
		if _, ok := palettes[v]; !ok {
			return errors.Errorf("unknown theme flavor %q", v)
		}

		o.flavor = v

		return nil
	}
}

// WithCapability pins the color tier instead of detecting it. Intended for
// tests: production code should let New detect the terminal.
func WithCapability(v Capability) Option {
	return func(o *Theme) error {
		o.capability = v
		return nil
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func New(opts ...Option) (*Theme, error) {
	inst := &Theme{
		flavor:     DefaultFlavor,
		capability: DetectCapability(),
	}

	for _, opt := range opts {
		if opt != nil {
			if err := opt(inst); err != nil {
				return nil, err
			}
		}
	}

	palette, ok := palettes[inst.flavor]
	if !ok {
		return nil, errors.Errorf("no palette for flavor %q", inst.flavor)
	}

	inst.hex = palette

	return inst, nil
}

// ------------------------------------------------------------------------------------------------
// ~ Getter
// ------------------------------------------------------------------------------------------------

func (t *Theme) Flavor() Flavor {
	return t.flavor
}

func (t *Theme) Capability() Capability {
	return t.capability
}

func (t *Theme) IsDark() bool {
	return t.flavor.IsDark()
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

// Hex returns the role's exact Catppuccin value.
//
// The prompt cannot render it — see Prompt — but it is the palette's source of
// truth and is exposed for commands that print their own colors.
func (t *Theme) Hex(r Role) string {
	return t.hex[r]
}

// Prompt returns the role as one of go-prompt's 16 ANSI colors.
//
// The prompt is capped at 16 colors, so this is an approximation. It is
// hand-mapped rather than derived: gookit's nearest-basic conversion collapses
// almost every Mocha role onto white, which would erase both the roles and the
// dark/light distinction in precisely the tier that cannot afford to lose them.
func (t *Theme) Prompt(r Role) goprompt.Color {
	if t.capability == CapabilityNone {
		return goprompt.DefaultColor
	}

	if m, ok := promptColors[t.flavor]; ok {
		if c, ok := m[r]; ok {
			return c
		}
	}

	return goprompt.DefaultColor
}

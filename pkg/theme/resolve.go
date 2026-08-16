package theme

import (
	"github.com/gookit/color"
)

// Resolve builds the theme for a configured flavor name.
//
// An empty name means no theme was configured and returns nil, which callers
// treat as "keep the previous, unthemed behavior" rather than as an error. An
// unrecognized name is an error: it is a typo in something the user wrote, and
// silently rendering a different theme would hide it.
func Resolve(flavor string) (*Theme, error) {
	if flavor == "" {
		return nil, nil //nolint:nilnil // no theme configured is not a failure
	}

	parsed, err := ParseFlavor(flavor)
	if err != nil {
		return nil, err
	}

	opts := []Option{WithFlavor(parsed)}
	if !color.SupportColor() {
		opts = append(opts, WithCapability(CapabilityNone))
	}

	return New(opts...)
}

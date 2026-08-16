package theme

import (
	"strings"

	"github.com/pkg/errors"
)

// Flavor is one of the four Catppuccin palettes.
type Flavor string

const (
	FlavorMocha     Flavor = "mocha"
	FlavorMacchiato Flavor = "macchiato"
	FlavorFrappe    Flavor = "frappe"
	FlavorLatte     Flavor = "latte"
)

// DefaultFlavor is used when nothing selects a flavor and auto-detection does
// not produce an answer.
const DefaultFlavor = FlavorMocha

// Flavors returns every flavor, dark ones first.
func Flavors() []Flavor {
	return []Flavor{FlavorMocha, FlavorMacchiato, FlavorFrappe, FlavorLatte}
}

// ParseFlavor resolves a flavor name case-insensitively.
func ParseFlavor(v string) (Flavor, error) {
	switch f := Flavor(strings.ToLower(strings.TrimSpace(v))); f {
	case FlavorMocha, FlavorMacchiato, FlavorFrappe, FlavorLatte:
		return f, nil
	default:
		return "", errors.Errorf("unknown theme flavor %q (expected one of: mocha, macchiato, frappe, latte)", v)
	}
}

func (f Flavor) String() string {
	return string(f)
}

// IsDark reports whether the flavor is intended for a dark terminal
// background. Only Latte is light.
func (f Flavor) IsDark() bool {
	return f != FlavorLatte
}

// palettes holds the hex value of every role, per flavor.
//
// Hex is the only hand-written color data in this package: the 256- and
// 16-color projections are computed from it at construction time, so the tiers
// cannot disagree with each other. Values are verified against upstream
// catppuccin/palette palette.json.
var palettes = map[Flavor]map[Role]string{
	FlavorMocha: {
		RoleText:      "#cdd6f4", // Text
		RoleMuted:     "#7f849c", // Overlay1
		RolePrimary:   "#cba6f7", // Mauve
		RoleSecondary: "#74c7ec", // Sapphire
		RoleSuccess:   "#a6e3a1", // Green
		RoleWarning:   "#f9e2af", // Yellow
		RoleError:     "#f38ba8", // Red
		RoleInfo:      "#89dceb", // Sky
		RoleDebug:     "#6c7086", // Overlay0
	},
	FlavorMacchiato: {
		RoleText:      "#cad3f5",
		RoleMuted:     "#8087a2",
		RolePrimary:   "#c6a0f6",
		RoleSecondary: "#7dc4e4",
		RoleSuccess:   "#a6da95",
		RoleWarning:   "#eed49f",
		RoleError:     "#ed8796",
		RoleInfo:      "#91d7e3",
		RoleDebug:     "#6e738d",
	},
	FlavorFrappe: {
		RoleText:      "#c6d0f5",
		RoleMuted:     "#838ba7",
		RolePrimary:   "#ca9ee6",
		RoleSecondary: "#85c1dc",
		RoleSuccess:   "#a6d189",
		RoleWarning:   "#e5c890",
		RoleError:     "#e78284",
		RoleInfo:      "#99d1db",
		RoleDebug:     "#737994",
	},
	FlavorLatte: {
		RoleText:      "#4c4f69",
		RoleMuted:     "#8c8fa1",
		RolePrimary:   "#8839ef",
		RoleSecondary: "#209fb5",
		RoleSuccess:   "#40a02b",
		RoleWarning:   "#df8e1d",
		RoleError:     "#d20f39",
		RoleInfo:      "#04a5e5",
		RoleDebug:     "#9ca0b0",
	},
}

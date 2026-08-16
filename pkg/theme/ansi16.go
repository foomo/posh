package theme

import (
	"github.com/foomo/posh/pkg/prompt/goprompt"
)

// The 16-color mapping is hand-written rather than computed.
//
// gookit's RGBColor.Basic() picks the nearest basic color by RGB distance,
// which for Catppuccin's low-saturation palette collapses almost every Mocha
// role onto white — nine roles rendered identically, and no visible difference
// between the dark and light flavors. Since go-prompt is permanently limited
// to these 16 slots, that is precisely where the distinctions matter most, so
// they are chosen by hand.
//
// Slots are picked to match Catppuccin's own published ANSI mapping
// (the `ansiColors` block of catppuccin/palette), so that on a terminal
// running a Catppuccin theme the prompt renders the intended colors exactly.
// That mapping puts each palette's signature colors in the *normal* slots
// (30-37) — Mocha's ANSI red is #f38ba8, the same value as its `red` — and
// uses the bright slots (90-97) for variant shades. Both dark and light
// flavors therefore use normal slots for the semantic roles; only the
// grays differ, since those must contrast with opposite backgrounds.
//
// The exceptions are deliberate:
//
//   - Primary maps to magenta even though Catppuccin's ANSI magenta is `pink`
//     rather than `mauve`. Magenta is the closest slot in hue, and picking by
//     RGB distance instead would land on white and erase the role.
//   - Secondary maps to blue and Info to cyan, matching `sapphire` and `sky`
//     by hue; both are near-neighbours of the exact ANSI values.
//
// The three dark flavors share a mapping: they differ only in shade, which 16
// colors cannot express anyway.

var promptDark = map[Role]goprompt.Color{
	RoleText:      goprompt.White,     // bright white: brightest against a dark bg
	RoleMuted:     goprompt.DarkGray,  // bright black
	RolePrimary:   goprompt.Purple,    // magenta
	RoleSecondary: goprompt.DarkBlue,  // blue
	RoleSuccess:   goprompt.DarkGreen, // green
	RoleWarning:   goprompt.Brown,     // yellow
	RoleError:     goprompt.DarkRed,   // red
	RoleInfo:      goprompt.Cyan,      // cyan
	RoleDebug:     goprompt.DarkGray,  // bright black
}

var promptLight = map[Role]goprompt.Color{
	RoleText:      goprompt.Black,     // darkest against a light bg
	RoleMuted:     goprompt.LightGray, // white: dimmed, but still legible
	RolePrimary:   goprompt.Purple,    // magenta
	RoleSecondary: goprompt.Cyan,      // cyan
	RoleSuccess:   goprompt.DarkGreen, // green
	RoleWarning:   goprompt.Brown,     // yellow
	RoleError:     goprompt.DarkRed,   // red
	RoleInfo:      goprompt.DarkBlue,  // blue
	RoleDebug:     goprompt.LightGray, // white
}

var promptColors = map[Flavor]map[Role]goprompt.Color{
	FlavorMocha:     promptDark,
	FlavorMacchiato: promptDark,
	FlavorFrappe:    promptDark,
	FlavorLatte:     promptLight,
}

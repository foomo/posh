package theme

import (
	"github.com/gookit/color"
)

// Capability is how much color the terminal can render.
type Capability int

const (
	CapabilityNone Capability = iota
	Capability16
	Capability256
	CapabilityTrueColor
)

// DetectCapability reports what the current terminal supports.
//
// This exists because neither pterm nor gookit degrades on its own: gookit's
// RenderCode only checks whether color is enabled at all, and otherwise emits
// whatever escape it was handed. A truecolor sequence therefore reaches a
// 256-color terminal verbatim and renders as garbage. Every color this package
// emits is chosen against this value instead.
func DetectCapability() Capability {
	switch {
	case !color.SupportColor():
		return CapabilityNone
	case color.SupportTrueColor():
		return CapabilityTrueColor
	case color.Support256Color():
		return Capability256
	default:
		return Capability16
	}
}

func (c Capability) String() string {
	switch c {
	case CapabilityNone:
		return "none"
	case Capability16:
		return "16"
	case Capability256:
		return "256"
	case CapabilityTrueColor:
		return "truecolor"
	default:
		return "unknown"
	}
}

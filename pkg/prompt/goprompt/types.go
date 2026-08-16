package goprompt

import (
	"github.com/c-bata/go-prompt"
)

type (
	Filter   = prompt.Filter
	Suggest  = prompt.Suggest
	Suggests = []prompt.Suggest
	Document = prompt.Document
	// Color is go-prompt's color type. It carries only the 16 ANSI slots below
	// — the underlying writer maps them through fixed SGR tables and has no
	// 256-color or truecolor path — so the interactive prompt cannot render
	// arbitrary hex.
	Color = prompt.Color
)

// The 16 ANSI slots go-prompt can render, low intensity first.
const (
	DefaultColor = prompt.DefaultColor
	Black        = prompt.Black
	DarkRed      = prompt.DarkRed
	DarkGreen    = prompt.DarkGreen
	Brown        = prompt.Brown
	DarkBlue     = prompt.DarkBlue
	Purple       = prompt.Purple
	Cyan         = prompt.Cyan
	LightGray    = prompt.LightGray
	DarkGray     = prompt.DarkGray
	Red          = prompt.Red
	Green        = prompt.Green
	Yellow       = prompt.Yellow
	Blue         = prompt.Blue
	Fuchsia      = prompt.Fuchsia
	Turquoise    = prompt.Turquoise
	White        = prompt.White
)

var (
	FilterFuzzy     = prompt.FilterFuzzy
	FilterContains  = prompt.FilterContains
	FilterHasPrefix = prompt.FilterHasPrefix
	FilterHasSuffix = prompt.FilterHasSuffix
	FilterCombined  = filterCombined
)

package prompt

import (
	"github.com/c-bata/go-prompt"
)

type (
	Filter   = prompt.Filter
	Suggest  = prompt.Suggest
	Document = prompt.Document
)

var (
	FilterFuzzy     = prompt.FilterFuzzy
	FilterContains  = prompt.FilterContains
	FilterHasPrefix = prompt.FilterHasPrefix
	FilterHasSuffix = prompt.FilterHasSuffix
)

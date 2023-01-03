package readline

import (
	"strings"
)

type Arg string

func (a Arg) String() string {
	return string(a)
}

func (a Arg) IsPipe() bool {
	return a == "|"
}

func (a Arg) IsPass() bool {
	return a == "--"
}

func (a Arg) IsFlag() bool {
	return strings.HasPrefix(a.String(), "-") && len(a) > 1
}

func (a Arg) IsRedirect() bool {
	return a == ">" || a == ">>" ||
		a == "2>" || a == "2>>" ||
		a == "&>" || a == "&>>" ||
		a == "2>&1"
}

func (a Arg) IsAdditional() bool {
	return a.IsPipe() || a.IsRedirect()
}

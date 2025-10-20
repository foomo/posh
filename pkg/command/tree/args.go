package tree

import (
	"github.com/foomo/posh/pkg/readline"
	"github.com/pkg/errors"
)

type Args []*Arg

func (a Args) Last() *Arg {
	if len(a) > 0 {
		return a[len(a)-1]
	} else {
		return nil
	}
}

func (a Args) Validate(args readline.Args) error {
	for j, arg := range a {
		if !arg.Optional && len(args)-1 < j {
			return errors.Wrap(ErrMissingArgument, arg.Name)
		}
	}

	return nil
}

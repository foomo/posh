package tree

import (
	"context"

	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
)

type Arg struct {
	Name     string
	Repeat   bool
	Optional bool
	Suggest  func(ctx context.Context, t *Root, r *readline.Readline) []goprompt.Suggest
}

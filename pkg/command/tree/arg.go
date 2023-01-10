package tree

import (
	"context"

	"github.com/c-bata/go-prompt"
	"github.com/foomo/posh/pkg/readline"
)

type Arg struct {
	Name     string
	Repeat   bool
	Optional bool
	Suggest  func(ctx context.Context, t *Root, r *readline.Readline) []prompt.Suggest
}

package command

import (
	"context"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/readline"
)

type Exit struct {
	l    log.Logger
	name string
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewExit(l log.Logger) *Exit {
	return &Exit{
		l:    l,
		name: "exit",
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (c *Exit) Name() string {
	return c.name
}

func (c *Exit) Description() string {
	return "exit shell"
}

// Skiller is deliberately not implemented: the command is inapplicable under
// `posh execute`, where each invocation is its own process and there is nothing
// to exit. A skill whose content is "do not use this" still has to be loaded
// before an agent can learn that, so it costs more than it saves - the root
// skill's index carries the name and description, which is all this warrants.

func (c *Exit) Execute(ctx context.Context, args *readline.Readline) error {
	return nil
}

func (c *Exit) Help(ctx context.Context, r *readline.Readline) string {
	return `Exit the Project Oriented Shell.

Usage:
  exit
`
}

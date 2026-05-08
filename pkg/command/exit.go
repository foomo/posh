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

func (c *Exit) Execute(ctx context.Context, args *readline.Readline) error {
	return nil
}

func (c *Exit) Help(ctx context.Context, r *readline.Readline) string {
	return `Exit the Project Oriented Shell.

Usage:
  exit
`
}

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

// Skill implements the optional Skiller interface. It exists to say the command
// is pointless outside the interactive shell, which its description implies to
// a human but not to an agent.
func (c *Exit) Skill(ctx context.Context) string {
	return "#### Notes\n\n" +
		"Leaves the interactive posh shell. There is nothing to exit under\n" +
		"`posh execute`, where each invocation is its own process, so running it\n" +
		"that way does nothing and succeeds - it is never the way to stop a\n" +
		"running command."
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

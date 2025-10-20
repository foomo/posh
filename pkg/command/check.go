package command

import (
	"context"

	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt/check"
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
)

type (
	Check struct {
		l        log.Logger
		tree     tree.Root
		check    check.Check
		checkers check.Checkers
	}
	Option func(*Check)
)

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewCheck(l log.Logger, checkers ...check.Checker) *Check {
	inst := &Check{
		l:        l,
		check:    check.DefaultCheck,
		checkers: checkers,
	}
	inst.tree = tree.New(&tree.Node{
		Name:        "check",
		Description: "Print all system checks",
		Execute:     inst.run,
	})

	return inst
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (c *Check) Name() string {
	return c.tree.Node().Name
}

func (c *Check) Description() string {
	return c.tree.Node().Description
}

func (c *Check) Complete(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
	return c.tree.Complete(ctx, r)
}

func (c *Check) Execute(ctx context.Context, r *readline.Readline) error {
	return c.tree.Execute(ctx, r)
}

func (c *Check) Help(ctx context.Context, r *readline.Readline) string {
	return c.tree.Help(ctx, r)
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (c *Check) run(ctx context.Context, r *readline.Readline) error {
	return c.check(ctx, c.l, c.checkers)
}

package command

import (
	"context"
	"strings"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt/history"
	"github.com/foomo/posh/pkg/readline"
)

type History struct {
	l       log.Logger
	name    string
	history history.History
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewHistory(l log.Logger, history history.History) *History {
	return &History{
		l:       l,
		name:    "history",
		history: history,
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (c *History) Name() string {
	return c.name
}

func (c *History) Description() string {
	return "show history"
}

func (c *History) Execute(ctx context.Context, r *readline.Readline) error {
	value, err := c.history.Load(ctx)
	if err != nil {
		return err
	}

	c.l.Info("History:\n\n" + strings.Join(value, "\n"))

	return nil
}

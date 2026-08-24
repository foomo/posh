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

// Skill implements the optional Skiller interface. What the one-line
// description cannot say is whose history this is, and that it is a record of
// intent rather than of success.
func (c *History) Skill(ctx context.Context, name string) string {
	return "#### Notes\n\n" +
		"Prints the commands previously entered in this project's interactive posh\n" +
		"shell, oldest last. It is normally file-backed (`.posh/history`) and so\n" +
		"survives across sessions, which makes it the best source for the real\n" +
		"invocations of a command - arguments and all - rather than guessing them\n" +
		"from the catalog.\n\n" +
		"Read it with two caveats. It records what was entered, not what\n" +
		"succeeded: a line may have failed, been a typo, or targeted something that\n" +
		"no longer exists. And repeated commands are collapsed to their most recent\n" +
		"occurrence, and the file is capped at a configured limit, so it is not a\n" +
		"chronological log and says nothing about how often a command is used.\n\n" +
		"Empty output is not an error - a project can configure history off, in\n" +
		"which case nothing is ever recorded."
}

func (c *History) Execute(ctx context.Context, r *readline.Readline) error {
	value, err := c.history.Load(ctx)
	if err != nil {
		return err
	}

	c.l.Info("History:\n\n" + strings.Join(value, "\n"))

	return nil
}

package command

import (
	"context"
	"math"
	"strings"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/pkg/errors"
)

type Help struct {
	l        log.Logger
	name     string
	commands Commands
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewHelp(l log.Logger, commands Commands) *Help {
	return &Help{
		l:        l,
		name:     "help",
		commands: commands,
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (c *Help) Name() string {
	return c.name
}

func (c *Help) Description() string {
	return "print help"
}

func (c *Help) Complete(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
	var suggests []goprompt.Suggest
	switch {
	case r.Args().LenLte(1):
		for _, value := range c.list() {
			suggests = append(suggests, goprompt.Suggest{Text: value.Name(), Description: value.Description()})
		}
	}
	return suggests
}

func (c *Help) Validate(ctx context.Context, r *readline.Readline) error {
	switch {
	case r.Args().LenIs(0):
		// all good
	case r.Args().LenIs(1):
		for _, command := range c.list() {
			if r.Args().At(0) == command.Name() {
				return nil
			}
		}
		return errors.Errorf("invalid [command] argument: %s", r.Args().At(0))
	}

	return nil
}

func (c *Help) Execute(ctx context.Context, r *readline.Readline) error {
	switch r.Args().Len() {
	case 0:
		ret := `Help about all available commands.

Usage:
  help [command]

Available Commands:
`
		for _, value := range c.list() {
			ret += c.format(value.Name(), value.Description())
		}
		c.l.Print(ret)
	default:
		if helper, ok := c.commands.Get(r.Args().At(0)).(Helper); ok {
			c.l.Print(helper.Help(ctx, r))
		} else {
			c.l.Print("command not found")
		}
	}
	return nil
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (c *Help) list() []Command {
	var ret []Command
	for _, value := range c.commands.List() {
		if _, ok := value.(Helper); ok {
			ret = append(ret, value)
		}
	}
	return ret
}

// print formatted output
func (c *Help) format(name, description string) string {
	offset := int(math.Max(0, float64(20-len(name))))
	suffix := strings.Repeat(" ", offset)
	prefix := ""
	if offset == 0 {
		suffix = "\n"
		prefix = strings.Repeat(" ", 20)
	}
	return "  " + name + suffix + prefix + description + "\n"
}

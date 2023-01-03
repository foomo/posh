package tree

import (
	"context"

	"github.com/c-bata/go-prompt"
	"github.com/foomo/posh/pkg/readline"

	"github.com/pkg/errors"
)

type Node struct {
	Name             string
	Names            func() []string
	Args             Args
	Flags            func(fs *readline.FlagSet)
	PassThroughArgs  Args
	PassThroughFlags func(fs *readline.FlagSet)
	Description      string
	Commands         []*Node
	Execute          func(ctx context.Context, args *readline.Readline) error
	//Suggest     func(ctx context.Context, parser *Parser, args *prompt.Args) (suggest []prompt.Suggest)
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (c *Node) setFlags(r *readline.Readline, parse bool) error {
	if c.Flags != nil {
		r.SetFlags(readline.NewFlagSet(c.Flags))
		if parse {
			if err := r.ParseFlags(); err != nil {
				return errors.Wrap(err, "failed to parse flags")
			}
		}
	}
	if c.PassThroughFlags != nil {
		r.SetParsePassThroughFlags(readline.NewFlagSet(c.PassThroughFlags))
		if parse {
			if err := r.ParsePassThroughFlags(); err != nil {
				return errors.Wrap(err, "failed to parse pass through flags")
			}
		}
	}
	return nil
}

func (c *Node) completeArguments(ctx context.Context, p *Root, r *readline.Readline, i int) (suggest []prompt.Suggest) {
	localArgs := r.Args()[i:]
	if len(c.Commands) > 0 && len(localArgs) <= 1 {
		for _, command := range c.Commands {
			suggest = append(suggest, prompt.Suggest{Text: command.Name, Description: command.Description})
		}
	} else if len(c.Args) >= len(localArgs) {
		if fn := c.Args[len(localArgs)-1].Suggest; fn != nil {
			suggest = fn(ctx, p, r)
		}
	} else if lastArg := c.Args.Last(); lastArg != nil && lastArg.Repeat {
		if fn := lastArg.Suggest; fn != nil {
			suggest = fn(ctx, p, r)
		}
	}
	return suggest
}

func (c *Node) completeFlags(r *readline.Readline) (suggest []prompt.Suggest) {
	for _, f := range r.AllFlags() {
		suggest = append(suggest, prompt.Suggest{Text: "--" + f.Name, Description: f.Usage})
	}
	return suggest
}

func (c *Node) completePassThroughFlags(r *readline.Readline) (suggest []prompt.Suggest) {
	for _, f := range r.AllPassThroughFlags() {
		suggest = append(suggest, prompt.Suggest{Text: "--" + f.Name, Description: f.Usage})
	}
	return suggest
}

func (c *Node) execute(ctx context.Context, r *readline.Readline, i int) error {
	localArgs := r.Args()[i:]
	if len(c.Commands) > 0 && len(localArgs) == 0 {
		return errors.New("missing [command] argument")
	} else if len(c.Args) > 0 {
		for j, arg := range c.Args {
			if !arg.Optional && len(localArgs) <= j+1 {
				return errors.New("missing [" + arg.Name + "] argument")
			}
		}
	}
	return c.Execute(ctx, r)
}

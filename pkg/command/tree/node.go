package tree

import (
	"context"

	"github.com/c-bata/go-prompt"
	"github.com/foomo/posh/pkg/readline"

	"github.com/pkg/errors"
)

type Node struct {
	Name             string
	Values           func(ctx context.Context, r *readline.Readline) []string
	Args             Args
	Flags            func(fs *readline.FlagSet)
	PassThroughArgs  Args
	PassThroughFlags func(fs *readline.FlagSet)
	Description      string
	Nodes            []*Node
	Execute          func(ctx context.Context, r *readline.Readline) error
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

func (c *Node) completeArguments(ctx context.Context, p *Root, r *readline.Readline, i int) []prompt.Suggest {
	var suggest []prompt.Suggest
	localArgs := r.Args()[i:]
	switch {
	case len(c.Nodes) > 0 && len(localArgs) <= 1:
		for _, command := range c.Nodes {
			if command.Values != nil {
				for _, name := range command.Values(ctx, r) {
					suggest = append(suggest, prompt.Suggest{Text: name, Description: command.Description})
				}
			} else {
				suggest = append(suggest, prompt.Suggest{Text: command.Name, Description: command.Description})
			}
		}
	case len(c.Args) > 0 && len(c.Args) >= len(localArgs):
		j := len(localArgs)
		if len(localArgs) > 0 && localArgs[j-1] != "" {
			j -= 1
		}
		if fn := c.Args[j].Suggest; fn != nil {
			suggest = fn(ctx, p, r)
		}
	case len(c.Args) > 0 && c.Args.Last().Repeat && c.Args.Last().Suggest != nil:
		suggest = c.Args.Last().Suggest(ctx, p, r)
	}
	return suggest
}

func (c *Node) completeFlags(r *readline.Readline) []prompt.Suggest {
	allFlags := r.AllFlags()
	suggest := make([]prompt.Suggest, len(allFlags))
	for i, f := range allFlags {
		suggest[i] = prompt.Suggest{Text: "--" + f.Name, Description: f.Usage}
	}
	return suggest
}

func (c *Node) completePassThroughFlags(r *readline.Readline) []prompt.Suggest {
	allPassThroughFlags := r.AllPassThroughFlags()
	suggest := make([]prompt.Suggest, len(allPassThroughFlags))
	for i, f := range allPassThroughFlags {
		suggest[i] = prompt.Suggest{Text: "--" + f.Name, Description: f.Usage}
	}
	return suggest
}

func (c *Node) execute(ctx context.Context, r *readline.Readline, i int) error {
	localArgs := r.Args()[i:]
	switch {
	case len(c.Nodes) > 0 && len(localArgs) == 0:
		return ErrMissingCommand
	case len(c.Args) > 0:
		for j, arg := range c.Args {
			if !arg.Optional && len(localArgs)-1 < j {
				return errors.Wrap(ErrMissingArgument, arg.Name)
			}
		}
	}
	return c.Execute(ctx, r)
}

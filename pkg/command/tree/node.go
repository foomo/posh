package tree

import (
	"context"
	"strings"

	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/foomo/posh/pkg/util/suggests"

	"github.com/pkg/errors"
)

type Node struct {
	Name             string
	Values           func(ctx context.Context, r *readline.Readline) []goprompt.Suggest
	Args             Args
	Flags            func(ctx context.Context, r *readline.Readline, fs *readline.FlagSet) error
	PassThroughArgs  Args
	PassThroughFlags func(ctx context.Context, r *readline.Readline, fs *readline.FlagSet) error
	Description      string
	Nodes            []*Node
	Execute          func(ctx context.Context, r *readline.Readline) error
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (c *Node) setFlags(ctx context.Context, r *readline.Readline, parse bool) error {
	if c.Flags != nil {
		f := readline.NewFlagSet()
		if err := c.Flags(ctx, r, f); err != nil {
			return err
		}
		r.SetFlags(f)
		if parse {
			if err := r.ParseFlags(); err != nil {
				return errors.Wrap(err, "failed to parse flags")
			}
		}
	}
	if c.PassThroughFlags != nil {
		f := readline.NewFlagSet()
		if err := c.PassThroughFlags(ctx, r, f); err != nil {
			return err
		}
		r.SetParsePassThroughFlags(f)
		if parse {
			if err := r.ParsePassThroughFlags(); err != nil {
				return errors.Wrap(err, "failed to parse pass through flags")
			}
		}
	}
	return nil
}

func (c *Node) completeArguments(ctx context.Context, p *Root, r *readline.Readline, i int) []goprompt.Suggest {
	var suggest []goprompt.Suggest
	localArgs := r.Args()[i:]
	switch {
	case len(c.Nodes) > 0 && len(localArgs) <= 1:
		for _, command := range c.Nodes {
			if command.Values != nil {
				suggest = command.Values(ctx, r)
			} else {
				suggest = append(suggest, goprompt.Suggest{Text: command.Name, Description: command.Description})
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

func (c *Node) completeFlags(r *readline.Readline) []goprompt.Suggest {
	allFlags := r.AllFlags()
	if r.Flags().LenGt(1) {
		if values := r.FlagSet().GetValues(strings.TrimPrefix(r.Flags().At(r.Flags().Len()-2), "--")); values != nil {
			return suggests.List(values)
		}
	}
	suggest := make([]goprompt.Suggest, len(allFlags))
	for i, f := range allFlags {
		suggest[i] = goprompt.Suggest{Text: "--" + f.Name, Description: f.Usage}
	}
	return suggest
}

func (c *Node) completePassThroughFlags(r *readline.Readline) []goprompt.Suggest {
	allPassThroughFlags := r.AllPassThroughFlags()
	suggest := make([]goprompt.Suggest, len(allPassThroughFlags))
	for i, f := range allPassThroughFlags {
		suggest[i] = goprompt.Suggest{Text: "--" + f.Name, Description: f.Usage}
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
	case c.Execute == nil && len(c.Nodes) > 0:
		return ErrMissingCommand
	case c.Execute == nil:
		return ErrInvalidCommand
	}
	return c.Execute(ctx, r)
}

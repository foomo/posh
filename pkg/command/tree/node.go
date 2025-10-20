package tree

import (
	"context"
	"strings"

	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
	xstrings "github.com/foomo/posh/pkg/util/strings"
	"github.com/foomo/posh/pkg/util/suggests"

	"github.com/pkg/errors"
)

type Node struct {
	Name        string
	Values      func(ctx context.Context, r *readline.Readline) []goprompt.Suggest
	Args        Args
	Flags       func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error
	Description string
	Nodes       []*Node
	Execute     func(ctx context.Context, r *readline.Readline) error
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (c *Node) setFlags(ctx context.Context, r *readline.Readline, parse bool) error {
	fs := readline.NewFlagSets()
	if c.Flags != nil {
		if err := c.Flags(ctx, r, fs); err != nil {
			return err
		}
	}

	r.SetFlagSets(fs)

	if parse {
		if err := r.ParseFlagSets(); err != nil {
			return errors.Wrap(err, "failed to parse flags")
		}
	}

	return nil
}

func (c *Node) completeArguments(ctx context.Context, p *root, r *readline.Readline, i int) []goprompt.Suggest {
	var suggest []goprompt.Suggest

	localArgs := r.Args().From(i)

	switch {
	case len(c.Nodes) > 0:
		for _, command := range c.Nodes {
			if command.Values != nil {
				suggest = command.Values(ctx, r)
			} else {
				suggest = append(suggest, goprompt.Suggest{Text: command.Name, Description: command.Description})
			}
		}
	case len(c.Args) > 0 && len(localArgs) == 0:
		if fn := c.Args[0].Suggest; fn != nil {
			suggest = fn(ctx, p, r)
		}
	case len(c.Args) > 0 && len(c.Args) >= len(localArgs):
		if fn := c.Args[len(localArgs)-1].Suggest; fn != nil {
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
		if fs := r.FlagSets(); fs != nil {
			if values := fs.All().GetValues(strings.TrimPrefix(r.Flags().At(r.Flags().Len()-2), "--")); values != nil {
				return suggests.List(values)
			}
		}
	}

	suggest := make([]goprompt.Suggest, len(allFlags))
	for i, f := range allFlags {
		suggest[i] = goprompt.Suggest{Text: "--" + f.Name, Description: f.Usage}
	}

	return suggest
}

func (c *Node) execute(ctx context.Context, r *readline.Readline, i int) error {
	localArgs := r.Args().From(i)

	// validate
	switch {
	case c.Execute == nil && len(c.Nodes) == 0:
		return ErrNoop
	case c.Execute == nil && len(c.Nodes) > 0 && len(localArgs) == 0:
		return ErrMissingCommand
	case c.Execute != nil && len(c.Nodes) == 0 && len(c.Args) == 0:
	case c.Execute != nil && len(c.Nodes) == 0 && c.Args.Validate(localArgs) == nil:
	case c.Execute != nil && len(c.Nodes) == 0 && c.Args.Validate(localArgs) != nil:
		return c.Args.Validate(localArgs)
	case c.Execute == nil:
		return ErrInvalidCommand
	}

	// execute
	return c.Execute(ctx, r)
}

func (c *Node) find(ctx context.Context, r *readline.Readline, i int) (*Node, int) {
	if r.Args().LenLt(i + 1) {
		return nil, i
	}

	arg := r.Args().At(i)
	for _, cmd := range c.Nodes {
		if cmd.Name == arg {
			if subCmd, j := cmd.find(ctx, r, i+1); subCmd != nil {
				return subCmd, j
			}

			return cmd, i + 1
		}

		if cmd.Values != nil {
			for _, name := range cmd.Values(ctx, r) {
				if name.Text == arg {
					if subCmd, j := cmd.find(ctx, r, i+1); subCmd != nil {
						return subCmd, j
					}

					return cmd, i + 1
				}
			}
		}
	}

	return nil, i
}

func (c *Node) help(ctx context.Context, r *readline.Readline) string {
	pad := "      "
	ret := c.Description

	if len(c.Nodes) > 0 {
		ret += "\n\nUsage:\n"
		ret += pad + c.Name + " [command]"

		ret += "\n\nAvailable Commands:\n"
		for _, node := range c.Nodes {
			ret += pad + xstrings.PadEnd(node.Name, " ", 30) + node.Description + "\n"
		}
	} else {
		ret += "\n\nUsage:\n"
		ret += pad + c.Name

		for _, arg := range c.Args {
			ret += " "
			if arg.Optional {
				ret += "<"
			} else {
				ret += "["
			}

			ret += arg.Name
			if arg.Optional {
				ret += ">"
			} else {
				ret += "]"
			}

			ret += "\n"
		}

		if len(c.Args) > 0 {
			ret += "\n\nArguments:\n"
			for _, arg := range c.Args {
				ret += pad + xstrings.PadEnd(arg.Name, " ", 30) + arg.Description + "\n"
			}
		}

		if c.Flags != nil {
			fs := readline.NewFlagSets()
			if err := c.Flags(ctx, r, fs); err == nil {
				ret += "\n\nFlags:\n"
				ret += fs.All().FlagUsages()
			}
		}
	}

	return ret
}

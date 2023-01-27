package tree

import (
	"context"
	"sort"

	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
)

type Root struct {
	Name        string
	Description string
	Node        *Node
	Nodes       Nodes
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (t *Root) Complete(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
	var suggests []goprompt.Suggest
	switch r.Mode() {
	case readline.ModeArgs:
		if r.Args().LenLte(1) && len(t.Nodes) > 0 {
			for _, command := range t.Nodes {
				if command.Values != nil {
					suggests = command.Values(ctx, r)
				} else {
					suggests = append(suggests, goprompt.Suggest{Text: command.Name, Description: command.Description})
				}
			}
		} else if cmd, i := t.find(ctx, t.Nodes, r, 0); cmd == nil && t.Node != nil {
			if err := t.Node.setFlags(ctx, r, false); err != nil {
				return nil
			} else {
				suggests = t.Node.completeArguments(ctx, t, r, 0)
			}
		} else if cmd == nil {
			return nil
		} else if err := cmd.setFlags(ctx, r, false); err != nil {
			return nil
		} else {
			suggests = cmd.completeArguments(ctx, t, r, i+1)
		}
	case readline.ModeFlags:
		if cmd, _ := t.find(ctx, t.Nodes, r, 0); cmd == nil && t.Node != nil {
			if err := t.Node.setFlags(ctx, r, false); err != nil {
				return nil
			} else {
				suggests = t.Node.completeFlags(r)
			}
		} else if cmd == nil {
			return nil
		} else if err := cmd.setFlags(ctx, r, false); err != nil {
			return nil
		} else {
			suggests = cmd.completeFlags(r)
		}
	case readline.ModePassThroughFlags:
		if cmd, _ := t.find(ctx, t.Nodes, r, 0); cmd == nil && t.Node != nil {
			if err := t.Node.setFlags(ctx, r, false); err != nil {
				return nil
			} else {
				suggests = t.Node.completePassThroughFlags(r)
			}
		} else if cmd == nil {
			return nil
		} else if err := cmd.setFlags(ctx, r, false); err != nil {
			return nil
		} else {
			suggests = cmd.completePassThroughFlags(r)
		}
	case readline.ModeAdditionalArgs:
		// do nothing
	}
	sort.Slice(suggests, func(i, j int) bool {
		return suggests[i].Text < suggests[j].Text
	})
	return suggests
}

func (t *Root) Execute(ctx context.Context, r *readline.Readline) error {
	var (
		cmd   *Node
		index int
	)

	switch {
	case t.Node == nil && len(t.Nodes) == 0:
		return ErrNoop
	case r.Args().LenIs(0) && t.Node == nil:
		return ErrMissingCommand
	}

	if r.Args().LenIs(0) {
		cmd = t.Node
	} else if found, i := t.find(ctx, t.Nodes, r, 0); found != nil {
		cmd = found
		index = i
	} else if t.Node == nil {
		return ErrInvalidCommand
	} else {
		cmd = t.Node
	}

	if err := cmd.setFlags(ctx, r, true); err != nil {
		return err
	} else if err := cmd.execute(ctx, r, index); err != nil {
		return err
	}
	return nil
}

func (t *Root) Help(ctx context.Context, r *readline.Readline) string {
	// TODO recursive help
	ret := t.Description
	if t.Nodes != nil {
		ret += "\n\nUsage:\n"
		ret += "  " + t.Name + " [command]"

		ret += "\n\nAvailable Commands:\n"
		for _, node := range t.Nodes {
			ret += "  " + node.Name
		}
	}
	return ret
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (t *Root) find(ctx context.Context, cmds []*Node, r *readline.Readline, i int) (*Node, int) {
	if r.Args().LenLt(i + 1) {
		return nil, i
	}
	arg := r.Args().At(i)
	for _, cmd := range cmds {
		if cmd.Name == arg {
			if subCmd, j := t.find(ctx, cmd.Nodes, r, i+1); subCmd != nil {
				return subCmd, j
			}
			return cmd, i
		}
		if cmd.Values != nil {
			for _, name := range cmd.Values(ctx, r) {
				if name.Text == arg {
					if subCmd, j := t.find(ctx, cmd.Nodes, r, i+1); subCmd != nil {
						return subCmd, j
					}
					return cmd, i
				}
			}
		}
	}
	return nil, i
}

package tree

import (
	"context"
	"sort"

	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
)

type Root interface {
	Node() *Node
	Complete(ctx context.Context, r *readline.Readline) []goprompt.Suggest
	Execute(ctx context.Context, r *readline.Readline) error
	Help(ctx context.Context, r *readline.Readline) string
}

type root struct {
	node *Node
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func New(node *Node) Root {
	return &root{
		node: node,
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (t *root) Node() *Node {
	return t.node
}

func (t *root) Complete(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
	var suggests []goprompt.Suggest
	switch r.Mode() {
	case readline.ModeArgs:
		if r.Args().LenLte(1) && len(t.node.Nodes) > 0 {
			for _, command := range t.node.Nodes {
				if command.Values != nil {
					suggests = command.Values(ctx, r)
				} else {
					suggests = append(suggests, goprompt.Suggest{Text: command.Name, Description: command.Description})
				}
			}
		} else if cmd, i := t.node.find(ctx, r, 0); cmd == nil && t.node != nil {
			if err := t.node.setFlags(ctx, r, false); err != nil {
				return nil
			} else {
				suggests = t.node.completeArguments(ctx, t, r, 0)
			}
		} else if cmd == nil {
			return nil
		} else if err := cmd.setFlags(ctx, r, false); err != nil {
			return nil
		} else {
			suggests = cmd.completeArguments(ctx, t, r, i+1)
		}
	case readline.ModeFlags:
		if cmd, _ := t.node.find(ctx, r, 0); cmd == nil && t.node != nil {
			if err := t.node.setFlags(ctx, r, false); err != nil {
				return nil
			} else {
				suggests = t.node.completeFlags(r)
			}
		} else if cmd == nil {
			return nil
		} else if err := cmd.setFlags(ctx, r, false); err != nil {
			return nil
		} else {
			suggests = cmd.completeFlags(r)
		}
	case readline.ModeAdditionalArgs:
		// do nothing
	}
	sort.Slice(suggests, func(i, j int) bool {
		return suggests[i].Text < suggests[j].Text
	})
	return suggests
}

func (t *root) Execute(ctx context.Context, r *readline.Readline) error {
	var (
		cmd   *Node
		index int
	)

	switch {
	case t.node == nil && t.node.Execute == nil && len(t.node.Nodes) == 0:
		return ErrNoop
	case r.Args().LenIs(1) && t.node == nil:
		return ErrMissingCommand
	}

	if r.Args().LenIs(0) {
		cmd = t.node
	} else if found, i := t.node.find(ctx, r, 0); found != nil {
		cmd = found
		index = i
	} else if t.node == nil {
		return ErrInvalidCommand
	} else {
		cmd = t.node
	}

	if err := cmd.setFlags(ctx, r, true); err != nil {
		return err
	} else if err := cmd.execute(ctx, r, index); err != nil {
		return err
	}
	return nil
}

func (t *root) Help(ctx context.Context, r *readline.Readline) string {
	var (
		cmd *Node
	)

	if t.node == nil {
		return "command not found"
	} else if r.Args().LenIs(1) {
		cmd = t.node
	} else if len(t.node.Nodes) == 0 {
		return "command not found"
	} else if found, _ := t.node.find(ctx, r, 1); found != nil {
		cmd = found
	} else {
		cmd = t.node
	}

	return cmd.help(ctx, r)
}

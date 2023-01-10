package tree

import (
	"context"
	"sort"

	"github.com/c-bata/go-prompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/pkg/errors"
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

func (t *Root) RunExecution(ctx context.Context, r *readline.Readline) error {
	var (
		cmd   *Node
		index int
	)
	if r.Args().LenIs(0) {
		cmd = t.Node
	} else if found, i := t.find(t.Nodes, r, 0); found != nil {
		cmd = found
		index = i
	} else {
		cmd = t.Node
	}

	if cmd == nil {
		return errors.New("invalid command")
	}

	if err := cmd.setFlags(r, true); err != nil {
		return err
	} else if err := cmd.execute(ctx, r, index); err != nil {
		return err
	}
	return nil
}

func (t *Root) RunCompletion(ctx context.Context, r *readline.Readline) []prompt.Suggest {
	var suggests []prompt.Suggest
	switch r.Mode() {
	case readline.ModeArgs:
		if r.Args().LenLte(1) && len(t.Nodes) > 0 {
			for _, command := range t.Nodes {
				suggests = append(suggests, prompt.Suggest{Text: command.Name, Description: command.Description})
			}
		} else if cmd, i := t.find(t.Nodes, r, 0); cmd == nil && t.Node != nil {
			if err := t.Node.setFlags(r, false); err != nil {
				return nil
			} else {
				suggests = t.Node.completeArguments(ctx, t, r, 0)
			}
		} else if cmd == nil {
			return nil
		} else if err := cmd.setFlags(r, false); err != nil {
			return nil
		} else {
			suggests = cmd.completeArguments(ctx, t, r, i+1)
		}
		sort.Slice(suggests, func(i, j int) bool {
			return suggests[i].Text < suggests[j].Text
		})
	case readline.ModeFlags:
		if cmd, _ := t.find(t.Nodes, r, 0); cmd == nil && t.Node != nil {
			if err := t.Node.setFlags(r, false); err != nil {
				return nil
			} else {
				suggests = t.Node.completeFlags(r)
			}
		} else if cmd == nil {
			return nil
		} else if err := cmd.setFlags(r, false); err != nil {
			return nil
		} else {
			suggests = cmd.completeFlags(r)
		}
		sort.Slice(suggests, func(i, j int) bool {
			return suggests[i].Text < suggests[j].Text
		})
	case readline.ModePassThroughArgs:
		// TODO
	case readline.ModePassThroughFlags:
		if cmd, _ := t.find(t.Nodes, r, 0); cmd == nil && t.Node != nil {
			if err := t.Node.setFlags(r, false); err != nil {
				return nil
			} else {
				suggests = t.Node.completePassThroughFlags(r)
			}
		} else if cmd == nil {
			return nil
		} else if err := cmd.setFlags(r, false); err != nil {
			return nil
		} else {
			suggests = cmd.completePassThroughFlags(r)
		}
		sort.Slice(suggests, func(i, j int) bool {
			return suggests[i].Text < suggests[j].Text
		})
	case readline.ModeAdditionalArgs:
		// do nothing
	}
	return suggests
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (t *Root) find(cmds []*Node, r *readline.Readline, i int) (*Node, int) {
	if r.Args().LenLt(i + 1) {
		return nil, i
	}
	arg := r.Args().At(i)
	for _, cmd := range cmds {
		if cmd.Name == arg {
			if subCmd, j := t.find(cmd.Nodes, r, i+1); subCmd != nil {
				return subCmd, j
			}
			return cmd, i
		}
		if cmd.Names != nil {
			for _, name := range cmd.Names() {
				if name == arg {
					if subCmd, j := t.find(cmd.Nodes, r, i+1); subCmd != nil {
						return subCmd, j
					}
				}
				return cmd, i
			}
		}
	}
	return nil, i
}

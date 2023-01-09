package tree

import (
	"context"
	"sort"

	"github.com/c-bata/go-prompt"
	"github.com/foomo/posh/pkg/readline"

	"github.com/pkg/errors"
)

type Root struct {
	Name  string
	Nodes []*Node
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (p *Root) Execute(ctx context.Context, r *readline.Readline) error {
	if r.Args().LenIs(0) {
		return errors.New("missing [command] argument")
	} else if cmd, i := p.find(p.Nodes, r, 0); cmd == nil {
		return errors.New("invalid [command] argument")
	} else if err := cmd.setFlags(r, true); err != nil {
		return err
	} else if err := cmd.execute(ctx, r, i); err != nil {
		return err
	}
	return nil
}

func (p *Root) Complete(ctx context.Context, r *readline.Readline) []prompt.Suggest {
	var suggests []prompt.Suggest
	switch r.Mode() {
	case readline.ModeArgs:
		if r.Args().LenLte(1) {
			for _, command := range p.Nodes {
				suggests = append(suggests, prompt.Suggest{Text: command.Name, Description: command.Description})
			}
		} else if cmd, i := p.find(p.Nodes, r, 0); cmd == nil {
			return nil
		} else if err := cmd.setFlags(r, false); err != nil {
			return nil
		} else {
			suggests = cmd.completeArguments(ctx, p, r, i+1)
		}
		sort.Slice(suggests, func(i, j int) bool {
			return suggests[i].Text < suggests[j].Text
		})
	case readline.ModeFlags:
		if cmd, _ := p.find(p.Nodes, r, 0); cmd == nil {
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
		if cmd, _ := p.find(p.Nodes, r, 0); cmd == nil {
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

func (p *Root) find(cmds []*Node, r *readline.Readline, i int) (*Node, int) {
	if r.Args().LenLt(i + 1) {
		return nil, i
	}
	arg := r.Args().At(i)
	for _, cmd := range cmds {
		if cmd.Name == arg {
			if subCmd, j := p.find(cmd.Nodes, r, i+1); subCmd != nil {
				return subCmd, j
			}
			return cmd, i
		}
		if cmd.Names != nil {
			for _, name := range cmd.Names() {
				if name == arg {
					if subCmd, j := p.find(cmd.Nodes, r, i+1); subCmd != nil {
						return subCmd, j
					}
				}
				return cmd, i
			}
		}
	}
	return nil, i
}

// Package complete produces shell-completion suggestions for `posh execute`
// invocations by reusing the same readline-driven dispatch that the
// interactive prompt uses.
package complete

import (
	"context"
	"strings"

	"github.com/foomo/posh/pkg/command"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
)

// Suggest parses (args, toComplete) into a readline state, dispatches to the
// matching command's completer interface, and returns suggestions in cobra's
// `value\tdescription` format. Suggestions are prefix-filtered by toComplete.
func Suggest(ctx context.Context, l log.Logger, cmds command.Commands, args []string, toComplete string) []string {
	r, err := readline.New(l)
	if err != nil {
		return nil
	}

	input := strings.Join(append(append([]string{}, args...), toComplete), " ")
	if err := r.Parse(input); err != nil {
		return nil
	}

	var suggests []goprompt.Suggest

	if r.IsModeDefault() && r.Args().LenIs(0) {
		for _, inst := range cmds.List() {
			suggests = append(suggests, goprompt.Suggest{Text: inst.Name(), Description: inst.Description()})
		}
	} else if cmd := cmds.Get(r.Cmd()); cmd != nil {
		switch r.Mode() {
		case readline.ModeArgs:
			if v, ok := cmd.(command.ArgumentCompleter); ok {
				suggests = v.CompleteArguments(ctx, r)
			} else if v, ok := cmd.(command.Completer); ok {
				suggests = v.Complete(ctx, r)
			}
		case readline.ModeFlags:
			if v, ok := cmd.(command.FlagCompleter); ok {
				suggests = v.CompleteFlags(ctx, r)
			} else if v, ok := cmd.(command.Completer); ok {
				suggests = v.Complete(ctx, r)
			}
		case readline.ModeAdditionalArgs:
			if v, ok := cmd.(command.AdditionalArgsCompleter); ok {
				suggests = v.CompleteAdditionalArgs(ctx, r)
			} else if v, ok := cmd.(command.Completer); ok {
				suggests = v.Complete(ctx, r)
			}
		}
	}

	out := make([]string, 0, len(suggests))
	for _, s := range suggests {
		if toComplete != "" && !strings.HasPrefix(s.Text, toComplete) {
			continue
		}

		if s.Description != "" {
			out = append(out, s.Text+"\t"+s.Description)
		} else {
			out = append(out, s.Text)
		}
	}

	return out
}

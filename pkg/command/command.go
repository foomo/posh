package command

import (
	"context"

	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
)

type (
	Command interface {
		Name() string
		Description() string
		Execute(ctx context.Context, r *readline.Readline) error
	}
	Helper interface {
		Help(ctx context.Context, r *readline.Readline) string
	}
	Validator interface {
		Validate(ctx context.Context, r *readline.Readline) error
	}
	Completer interface {
		Complete(ctx context.Context, r *readline.Readline) []goprompt.Suggest
	}
	FlagCompleter interface {
		CompleteFlags(ctx context.Context, r *readline.Readline) []goprompt.Suggest
	}
	ArgumentCompleter interface {
		CompleteArguments(ctx context.Context, r *readline.Readline) []goprompt.Suggest
	}
	PassThroughFlagsCompleter interface {
		CompletePassTroughFlags(ctx context.Context, r *readline.Readline) []goprompt.Suggest
	}
	AdditionalArgsCompleter interface {
		CompleteAdditionalArgs(ctx context.Context, r *readline.Readline) []goprompt.Suggest
	}
)

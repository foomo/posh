package command

import (
	"context"

	"github.com/c-bata/go-prompt"
	"github.com/foomo/posh/pkg/readline"
)

type (
	Command interface {
		Name() string
		Description() string
		Execute(ctx context.Context, r *readline.Readline) error
	}
	Helper interface {
		Help() string
	}
	Validator interface {
		Validate(ctx context.Context, r *readline.Readline) error
	}
	Completer interface {
		Complete(ctx context.Context, r *readline.Readline, d prompt.Document) []prompt.Suggest
	}
	FlagCompleter interface {
		CompleteFlags(ctx context.Context, r *readline.Readline, d prompt.Document) []prompt.Suggest
	}
	ArgumentCompleter interface {
		CompleteArguments(ctx context.Context, r *readline.Readline, d prompt.Document) []prompt.Suggest
	}
	PassThroughArgsCompleter interface {
		CompletePassTroughArgs(ctx context.Context, r *readline.Readline, d prompt.Document) []prompt.Suggest
	}
	PassThroughFlagsCompleter interface {
		CompletePassTroughFlags(ctx context.Context, r *readline.Readline, d prompt.Document) []prompt.Suggest
	}
	AdditionalArgsCompleter interface {
		CompleteAdditionalArgs(ctx context.Context, r *readline.Readline, d prompt.Document) []prompt.Suggest
	}
)

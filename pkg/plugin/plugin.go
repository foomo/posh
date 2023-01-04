package plugin

import (
	"context"

	"github.com/foomo/posh/pkg/config"
)

type Plugin interface {
	Prompt(ctx context.Context, cfg config.Prompt) error
	Execute(ctx context.Context, args []string) error
	Packages(ctx context.Context, cfg []config.Package) error
	Dependencies(ctx context.Context, cfg config.Dependencies) error
}

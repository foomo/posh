package plugin

import (
	"context"

	"github.com/foomo/posh/pkg/config"
)

type Plugin interface {
	Prompt(ctx context.Context, cfg config.Prompt) error
	Execute(ctx context.Context, args []string) error
	Brew(ctx context.Context, cfg config.Ownbrew) error
	Require(ctx context.Context, cfg config.Require) error
}

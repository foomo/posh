package plugin

import (
	"context"

	ownbrewconfig "github.com/foomo/ownbrew/pkg/config"
	"github.com/foomo/posh/pkg/config"
)

type Plugin interface {
	Prompt(ctx context.Context, cfg config.Prompt) error
	Execute(ctx context.Context, args []string) error
	Brew(ctx context.Context, cfg ownbrewconfig.Config, tags []string, dry bool) error
	Require(ctx context.Context, cfg config.Require) error
}

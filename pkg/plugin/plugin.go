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

// Completer is an optional Plugin extension that produces shell completion
// suggestions for `posh execute`. Returned strings use the cobra format:
// "value\tdescription" (description optional).
type Completer interface {
	Complete(ctx context.Context, args []string, toComplete string) []string
}

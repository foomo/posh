package exec

import (
	"context"
)

type CommandProvider func(ctx context.Context, args ...string) *Command

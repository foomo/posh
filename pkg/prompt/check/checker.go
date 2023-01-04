package check

import (
	"context"

	"github.com/foomo/posh/pkg/log"
)

type Checker func(ctx context.Context, l log.Logger) Info

package check

import (
	"context"

	"github.com/foomo/posh/pkg/log"
)

type Check func(ctx context.Context, l log.Logger, checkers []Checker) error

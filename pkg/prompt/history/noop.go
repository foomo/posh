package history

import (
	"context"

	"github.com/foomo/posh/pkg/log"
)

type Noop struct {
	l log.Logger
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewNoop(l log.Logger) *Noop {
	return &Noop{
		l: l,
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (h *Noop) Load(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (h *Noop) Persist(ctx context.Context, record string) {
}

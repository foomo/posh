package history

import (
	"context"
)

type History interface {
	Load(ctx context.Context) ([]string, error)
	Persist(ctx context.Context, record string)
}

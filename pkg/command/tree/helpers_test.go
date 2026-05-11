package tree_test

import (
	"context"
	"testing"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/readline"
	"github.com/stretchr/testify/require"
)

type ctxKey struct{}

func T(ctx context.Context) *testing.T {
	return ctx.Value(ctxKey{}).(*testing.T)
}

func SetT(ctx context.Context, t *testing.T) context.Context {
	t.Helper()
	return context.WithValue(ctx, ctxKey{}, t)
}

func newReadline(t *testing.T, input string) *readline.Readline {
	t.Helper()
	rl, err := readline.New(log.NewTest(t, log.TestWithLevel(log.LevelInfo)))
	require.NoError(t, err)
	require.NoError(t, rl.Parse(input))

	return rl
}

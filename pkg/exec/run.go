package exec

import (
	"context"
	"os/exec"
	"slices"
)

func Run(ctx context.Context, cmd *exec.Cmd, middlewares ...Middleware) error {
	run := func(_ context.Context, cmd *exec.Cmd) error {
		return cmd.Run()
	}

	for _, v := range slices.Backward(middlewares) {
		run = v(run)
	}

	return run(ctx, cmd)
}

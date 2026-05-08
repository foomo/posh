package exec

import (
	"context"
	"os/exec"
)

type Handler func(ctx context.Context, cmd *exec.Cmd) error

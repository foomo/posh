package middleware

import (
	"bytes"
	"context"
	"os/exec"

	pkgexec "github.com/foomo/posh/pkg/exec"
)

func WithEnv(v ...string) pkgexec.Middleware {
	return func(next pkgexec.Handler) pkgexec.Handler {
		return func(ctx context.Context, cmd *exec.Cmd) error {
			cmd.Env = append(cmd.Env, v...)
			return next(ctx, cmd)
		}
	}
}

func CaptureStdout(buf *bytes.Buffer) pkgexec.Middleware {
	return func(next pkgexec.Handler) pkgexec.Handler {
		return func(ctx context.Context, cmd *exec.Cmd) error {
			cmd.Stdout = buf
			return next(ctx, cmd)
		}
	}
}

func CaptureStderr(buf *bytes.Buffer) pkgexec.Middleware {
	return func(next pkgexec.Handler) pkgexec.Handler {
		return func(ctx context.Context, cmd *exec.Cmd) error {
			cmd.Stderr = buf
			return next(ctx, cmd)
		}
	}
}

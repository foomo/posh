package exec

import (
	"context"
	"io"
	"os"
	"os/exec"
)

type (
	Command struct {
		ctx         context.Context
		cmd         *exec.Cmd
		middlewares []Middleware
	}
)

func NewCommand(ctx context.Context, name string, arg ...string) *Command {
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return &Command{
		ctx: ctx,
		cmd: cmd,
	}
}

func (c *Command) Args(args ...string) *Command {
	c.cmd.Args = append(c.cmd.Args, args...)
	return c
}

func (c *Command) Env(env ...string) *Command {
	c.cmd.Env = append(c.cmd.Env, env...)
	return c
}

func (c *Command) Dir(dir string) *Command {
	c.cmd.Dir = dir
	return c
}

func (c *Command) Stdin(v io.Reader) *Command {
	c.cmd.Stdin = v
	return c
}

func (c *Command) Stdout(v io.Writer) *Command {
	c.cmd.Stdout = v
	return c
}

func (c *Command) Stderr(v io.Writer) *Command {
	c.cmd.Stderr = v
	return c
}

func (c *Command) Middleware(mw ...Middleware) *Command {
	c.middlewares = append(c.middlewares, mw...)
	return c
}

func (c *Command) Run() error {
	return Run(c.ctx, c.cmd, c.middlewares...)
}

package exec

import (
	"context"
	"io"
)

type stdioKey struct{}

type stdio struct {
	stdout io.Writer
	stderr io.Writer
}

// WithStdio returns a context that overrides the default Stdout/Stderr for
// any Command or Shell constructed with it, instead of the real process
// stdio. Used by posh mcp's RunCommand so subprocess output a command
// shells out to lands in the same buffer as its own logging, rather than
// corrupting the MCP stdio transport.
func WithStdio(ctx context.Context, stdout, stderr io.Writer) context.Context {
	return context.WithValue(ctx, stdioKey{}, stdio{stdout: stdout, stderr: stderr})
}

// StdioFrom returns the stdout/stderr writers set by WithStdio, if any.
func StdioFrom(ctx context.Context) (stdout, stderr io.Writer, ok bool) {
	v, ok := ctx.Value(stdioKey{}).(stdio)
	return v.stdout, v.stderr, ok
}

package mcp

import (
	"bytes"
	"context"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/pkg/errors"
)

// RunCommand resolves a fresh plugin instance via provider, wired to a logger
// that buffers instead of writing to stdout, executes args against it, and
// returns everything the command logged. The error, if any, is returned
// alongside the buffered output rather than swallowing it, since a failing
// command's own log lines are often the only diagnostic.
func RunCommand(ctx context.Context, provider plugin.Provider, args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("missing [cmd] argument")
	}

	var buf bytes.Buffer

	l := log.NewAgentJSON(log.AgentJSONWithWriter(&buf))

	plg, err := provider(l)
	if err != nil {
		return buf.String(), err
	}

	if err := plg.Execute(ctx, args); err != nil {
		return buf.String(), err
	}

	return buf.String(), nil
}

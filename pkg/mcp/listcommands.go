package mcp

import (
	"context"

	"github.com/foomo/posh/pkg/plugin"
)

// ListCommands returns the same catalog posh agent catalog prints, using the
// same plugin.List fallback plg already supports.
//
// path and depth scope the result for projects whose catalog is too large to
// hand an agent in full every call:
//
//   - path, if non-empty, keeps only the top-level command whose FullPath
//     matches it exactly (e.g. "squadron"), dropping every other top-level
//     command. A path that matches nothing returns an empty catalog rather
//     than an error - the caller asked for a subtree that doesn't exist.
//   - depth, if greater than 0, truncates Subcommands beyond that many
//     levels below whatever commands remain after the path filter. depth <= 0
//     means unlimited (the full tree, matching prior behavior).
func ListCommands(ctx context.Context, plg any, path string, depth int) (plugin.Catalog, error) {
	commands, err := plugin.List(ctx, plg)
	if err != nil {
		return plugin.Catalog{}, err
	}

	if path != "" {
		commands = filterPath(commands, path)
	}

	if depth > 0 {
		commands = truncateDepth(commands, depth)
	}

	return plugin.NewCatalog(commands), nil
}

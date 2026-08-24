package mcp

import (
	"context"

	"github.com/foomo/posh/pkg/plugin"
)

// ListCommands returns the same catalog posh agent catalog prints, using the
// same plugin.List fallback plg already supports.
func ListCommands(ctx context.Context, plg any) (plugin.Catalog, error) {
	commands, err := plugin.List(ctx, plg)
	if err != nil {
		return plugin.Catalog{}, err
	}

	return plugin.NewCatalog(commands), nil
}

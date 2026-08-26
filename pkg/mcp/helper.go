package mcp

import (
	"github.com/foomo/posh/pkg/plugin"
)

// filterPath keeps only the top-level command whose FullPath equals path.
func filterPath(commands []plugin.CommandInfo, path string) []plugin.CommandInfo {
	for _, c := range commands {
		if c.FullPath == path {
			return []plugin.CommandInfo{c}
		}
	}

	return nil
}

// truncateDepth returns a copy of commands with Subcommands cleared depth
// levels below them, so a caller can bound how much of a deep tree comes back
// without losing every command below the cutoff (they are just no longer
// expanded). depth counts levels below commands itself: depth 1 keeps
// commands' immediate Subcommands but clears their grandchildren.
func truncateDepth(commands []plugin.CommandInfo, depth int) []plugin.CommandInfo {
	ret := make([]plugin.CommandInfo, len(commands))

	for i, c := range commands {
		ret[i] = c

		if depth <= 0 {
			ret[i].Subcommands = nil
		} else {
			ret[i].Subcommands = truncateDepth(c.Subcommands, depth-1)
		}
	}

	return ret
}

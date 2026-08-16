package command

import (
	"context"

	"github.com/foomo/posh/pkg/command/tree"
)

// Describer is an optional Command extension: a command that describes its own
// structure for `posh agent catalog`, so an AI coding agent gets a
// machine-readable answer to "what can I run" without parsing help text.
//
// The contract is CommandInfo, not a command tree: a command is free to model
// its structure however it likes and stays free to change that later. Tree
// based commands delegate to tree.Root.Describe like any other tree method:
//
//	func (c *Cache) Describe(ctx context.Context) command.CommandInfo {
//		return c.tree.Describe(ctx)
//	}
//
// Commands that do not implement it are still listed, just without subcommand
// detail.
type Describer interface {
	Describe(ctx context.Context) CommandInfo
}

// The catalog shapes are built by the tree walk and aliased here so a command
// describing itself by hand does not need to import the tree package.
type (
	// CommandInfo describes one command for the `posh agent catalog` catalog.
	CommandInfo = tree.CommandInfo
	// ArgInfo describes a single positional argument of a command.
	ArgInfo = tree.ArgInfo
	// FlagInfo describes a single flag of a command.
	FlagInfo = tree.FlagInfo
)

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

// Describe builds a CommandInfo for a single command. Commands implementing
// the optional Describer interface describe themselves; anything else becomes a
// leaf carrying only its name and description.
func Describe(ctx context.Context, name, description string, v any) CommandInfo {
	if d, ok := v.(Describer); ok {
		return d.Describe(ctx)
	}

	return CommandInfo{FullPath: name, Description: description}
}

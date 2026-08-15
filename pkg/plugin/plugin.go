package plugin

import (
	"context"

	ownbrewconfig "github.com/foomo/ownbrew/pkg/config"
	"github.com/foomo/posh/pkg/command"
	"github.com/foomo/posh/pkg/config"
)

type Plugin interface {
	Prompt(ctx context.Context, cfg config.Prompt) error
	Execute(ctx context.Context, args []string) error
	Brew(ctx context.Context, cfg ownbrewconfig.Config, tags []string, dry bool) error
	Require(ctx context.Context, cfg config.Require) error
}

// Completer is an optional Plugin extension that produces shell completion
// suggestions for `posh execute`. Returned strings use the cobra format:
// "value\tdescription" (description optional).
type Completer interface {
	Complete(ctx context.Context, args []string, toComplete string) []string
}

// Lister is an optional Plugin extension that exposes the catalog of
// commands, used by `posh agent catalog` to give an AI coding agent a
// machine-readable answer to "what can I run" without parsing human-formatted
// help text.
type Lister interface {
	List(ctx context.Context) []CommandInfo
}

// SkillLister is an optional Plugin extension that lists commands together with
// the extra SKILL.md markdown each one contributes via command.Skiller.
//
// A plugin implementing only Lister still renders a skill, just without the
// extra prose.
type SkillLister interface {
	ListSkill(ctx context.Context) []SkillCommand
}

// SkillMetadataer is an optional Plugin extension supplying the generated
// SKILL.md frontmatter. Unset fields fall back to posh's defaults.
type SkillMetadataer interface {
	SkillMetadata(ctx context.Context) SkillMetadata
}

// Catalog is the payload of `posh agent catalog`.
//
// Unlike command results, which are emitted bare, the catalog keeps its
// type/schema_version discriminator: the shape is a published contract mirroring
// other agent-facing CLIs, and consumers dispatch on the marker.
type Catalog struct {
	Type          string        `json:"type"`
	SchemaVersion string        `json:"schema_version"`
	Commands      []CommandInfo `json:"commands"`
}

// CatalogType is the discriminator carried by Catalog.
const CatalogType = "posh.command_catalog"

// CatalogSchemaVersion is the current version of the Catalog shape.
const CatalogSchemaVersion = "3"

// NewCatalog wraps commands in a Catalog with the current type and version.
func NewCatalog(commands []CommandInfo) Catalog {
	return Catalog{
		Type:          CatalogType,
		SchemaVersion: CatalogSchemaVersion,
		Commands:      commands,
	}
}

// The catalog shapes live in pkg/command, next to the Command interface they
// describe, and are aliased here so a plugin does not need to import both.
type (
	// CommandInfo describes one command for the `posh agent catalog` catalog.
	CommandInfo = command.CommandInfo
	// ArgInfo describes a single positional argument of a command.
	ArgInfo = command.ArgInfo
	// FlagInfo describes a single flag of a command.
	FlagInfo = command.FlagInfo
)

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

// Describe builds a CommandInfo for a single command, see command.Describe.
func Describe(ctx context.Context, name, description string, v any) CommandInfo {
	return command.Describe(ctx, name, description, v)
}

// Skill returns the extra SKILL.md markdown a command contributes, see
// command.Skill.
func Skill(ctx context.Context, v any) string {
	return command.Skill(ctx, v)
}

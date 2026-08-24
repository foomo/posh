package command

import (
	"context"
)

// Skiller is an optional Command extension: a command that contributes extra
// markdown to its generated SKILL.md, beyond the structure Describer already
// provides - extended instructions, worked configuration examples, links to
// further docs.
//
// The contract is deliberately a free-form string rather than a struct: what is
// worth telling an agent varies per command, and a fixed shape would constrain
// that without buying anything, since the destination is markdown either way.
//
// name is the command name the contribution is rendered under, which is not
// necessarily the command's default one: the same provider can be registered
// under another name. Write the prose against name rather than hardcoding a
// literal, or the generated skill tells an agent to run a command that does not
// exist in this project.
//
// The returned markdown is appended verbatim under the command's heading. Use
// heading level 4 (####) or deeper - level 3 is the command heading itself.
//
// Commands that do not implement it render as just their paths and descriptions.
type Skiller interface {
	Skill(ctx context.Context, name string) string
}

// SkillMetadataer is an optional Command extension supplying the frontmatter of
// that command's own generated skill.
//
// Unset fields fall back to values derived from the command: the name to
// "posh-<name>", the description to a generic "Use when running ..." line. The
// description is the only thing an agent runtime sees when deciding whether to
// load a skill at all, so a generic one effectively means the skill never
// triggers - name the conditions that should reach for this command.
//
// name is the registered command name, as for Skiller.
type SkillMetadataer interface {
	SkillMetadata(ctx context.Context, name string) SkillMetadata
}

// SkillMetadata is the frontmatter of a generated SKILL.md.
//
// The fields are restricted to the Agent Skills spec allowlist
// (https://agentskills.io) so a generated skill stays valid if it is ever
// uploaded to claude.ai or the Skills API, both of which reject unknown keys.
type SkillMetadata struct {
	// Name defaults to "posh" for the root skill, "posh-<command>" for a
	// command's own.
	Name string `yaml:"name"`
	// Description defaults to a generic one. A description naming concrete
	// triggering conditions is what makes the skill load at all.
	Description string `yaml:"description"`
	// AllowedTools are granted without a permission prompt for the turn that
	// invokes the skill, e.g. "Bash(posh execute:*)".
	AllowedTools []string `yaml:"allowed-tools,omitempty"`
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

// Skill returns the extra SKILL.md markdown v contributes under name, or an
// empty string if it does not implement Skiller.
func Skill(ctx context.Context, v any, name string) string {
	if s, ok := v.(Skiller); ok {
		return s.Skill(ctx, name)
	}

	return ""
}

// SkillMetadataOf returns the skill frontmatter v supplies for name, or the zero
// value - which renders the derived defaults - when it does not implement
// SkillMetadataer.
func SkillMetadataOf(ctx context.Context, v any, name string) SkillMetadata {
	if s, ok := v.(SkillMetadataer); ok {
		return s.SkillMetadata(ctx, name)
	}

	return SkillMetadata{}
}

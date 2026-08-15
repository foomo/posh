package command

import (
	"context"
)

// Skiller is an optional Command extension: a command that contributes extra
// markdown to the generated SKILL.md, beyond the structure Describer already
// provides - extended instructions, worked configuration examples, links to
// further docs.
//
// The contract is deliberately a free-form string rather than a struct: what is
// worth telling an agent varies per command, and a fixed shape would constrain
// that without buying anything, since the destination is markdown either way.
//
// The returned markdown is appended verbatim under the command's heading in
// SKILL.md, after its arguments, flags and subcommands. Use heading level 4
// (####) or deeper - level 3 is the command heading itself.
//
// Commands that do not implement it render exactly as they did before.
type Skiller interface {
	Skill(ctx context.Context) string
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

// Skill returns the extra SKILL.md markdown v contributes, or an empty string if
// it does not implement Skiller.
func Skill(ctx context.Context, v any) string {
	if s, ok := v.(Skiller); ok {
		return s.Skill(ctx)
	}

	return ""
}

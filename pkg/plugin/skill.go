package plugin

import (
	"fmt"
	"strings"

	"github.com/foomo/posh/pkg/util"
	"gopkg.in/yaml.v3"
)

// DefaultSkillPath is where `posh agent skill install` writes the skill unless
// given an explicit path.
const DefaultSkillPath = ".claude/skills/posh/SKILL.md"

// Defaults for the generated frontmatter, used for any SkillMetadata field a
// project leaves unset.
const (
	defaultSkillName        = "posh"
	defaultSkillDescription = "Drive this project's posh shell (Project Oriented Shell)."
)

// skillSetup documents posh's own subcommands, as opposed to the project
// commands the catalog below is generated from.
//
// It is a fixed string rather than a walk of the cobra tree: these commands are
// the same in every project, they are wired up in cmd/ - which imports this
// package, so the dependency cannot go the other way - and an agent mostly
// needs them when `posh execute` fails because the environment is not set up.
const skillSetup = "## Setup\n\n" +
	"All of it is configured in this project's `.posh.yaml`; read that file to see\n" +
	"what is actually required or installed before running anything below.\n\n" +
	"- `posh require` - validate the preconditions from the `require` key (env\n" +
	"  vars, packages, scripts). Run it first when a command fails unexpectedly.\n" +
	"- `posh brew` - install the pinned tool versions from the `ownbrew` key. Run\n" +
	"  it when `posh require` reports a missing package; `--dry` only prints what\n" +
	"  would be installed.\n" +
	"- `posh prompt` - the interactive shell, configured by the `prompt` key. Do\n" +
	"  not run it under an agent harness: it blocks waiting on a TTY. Use\n" +
	"  `posh execute` for individual commands instead.\n\n"

// SkillCommand pairs a described command with the extra markdown its
// command.Skiller implementation contributed, if any.
//
// The catalog shape stays untouched: skill prose is a rendering concern, so it
// rides alongside CommandInfo rather than inside it.
type SkillCommand struct {
	CommandInfo

	// Skill is the markdown contributed by command.Skiller, or empty.
	Skill string
}

// SkillMetadata is the SKILL.md frontmatter a project can override.
//
// The fields are restricted to the Agent Skills spec allowlist
// (https://agentskills.io) so the generated skill stays valid if it is ever
// uploaded to claude.ai or the Skills API, both of which reject unknown keys.
type SkillMetadata struct {
	// Name defaults to "posh".
	Name string `yaml:"name"`
	// Description defaults to a generic one. A project specific description is
	// what makes the skill trigger on this project's vocabulary.
	Description string `yaml:"description"`
	// AllowedTools are granted without a permission prompt for the turn that
	// invokes the skill, e.g. "Bash(posh execute:*)".
	AllowedTools []string `yaml:"allowed-tools,omitempty"`
}

// RenderSkill renders a Claude Code SKILL.md from a command catalog.
//
// The result is generated in full from the catalog every time, so it always
// reflects whatever commands a project's posh shell currently registers - there
// is no hand-maintained content to drift out of sync. A zero SkillMetadata
// renders posh's default frontmatter.
func RenderSkill(meta SkillMetadata, commands []SkillCommand) string {
	var b strings.Builder

	writeSkillFrontmatter(&b, meta)

	b.WriteString("# posh\n\n")
	b.WriteString("Run project tasks via `posh execute <command...>` (alias `posh x`) instead of\n")
	b.WriteString("guessing shell commands directly. JSON output is automatic under an agent\n")
	b.WriteString("harness; `POSH_AGENT_MODE=0` forces human output. Re-run\n")
	b.WriteString("`posh agent skill update` after this project's commands change to refresh\n")
	b.WriteString("this file.\n\n")
	b.WriteString("posh does not enforce access control: whatever the surrounding harness or CI\n")
	b.WriteString("permits is the real boundary. A command being listed here does not mean it is\n")
	b.WriteString("safe to run unattended.\n\n")
	b.WriteString(skillSetup)
	b.WriteString("## Commands\n\n")
	b.WriteString("Only commands with more to say than their one-line description appear below.\n")
	b.WriteString("This is not the full list - run `posh agent catalog` for every command as\n")
	b.WriteString("JSON, or `posh help` for the human-readable list.\n\n")

	for _, c := range commands {
		if !c.describes() {
			continue
		}

		heading := c.FullPath
		if usage := c.Usage(); usage != "" {
			heading += " " + usage
		}

		fmt.Fprintf(&b, "### `%s`\n\n%s\n\n", heading, c.Description)

		// Arguments, flags and subcommands share one heading: without it the
		// list runs straight on from the description, while the command's own
		// Skiller prose below does carry a heading - so the structure reads as
		// part of the description.
		if len(c.Arguments) > 0 || len(c.Flags) > 0 || len(c.Subcommands) > 0 {
			b.WriteString("#### Usage\n\n")
		}

		writeSkillDetails(&b, c.CommandInfo, "")

		if len(c.Arguments) > 0 || len(c.Flags) > 0 {
			b.WriteString("\n")
		}

		writeSkillSubcommands(&b, c.Subcommands, 0)

		// Trimmed so a contribution ending in a newline does not stack blank
		// lines on top of the one added here.
		if skill := strings.TrimSpace(c.Skill); skill != "" {
			b.WriteString(skill + "\n\n")
		}
	}

	return b.String()
}

// describes reports whether a command says anything the skill is worth spending
// tokens on: prose from command.Skiller, or structure from command.Describer.
//
// The test is deliberately structural rather than "does it implement the
// interfaces". command.Describe already collapses a non-Describer into a
// CommandInfo carrying only a path and description, so by the time it arrives
// here a command that never opted in is indistinguishable from a Describer that
// legitimately describes a bare leaf - and both are equally uninformative. An
// agent reading the catalog has the name and description either way.
func (c SkillCommand) describes() bool {
	return strings.TrimSpace(c.Skill) != "" ||
		len(c.Arguments) > 0 ||
		len(c.Flags) > 0 ||
		len(c.Subcommands) > 0
}

// writeSkillFrontmatter renders the YAML frontmatter, falling back to posh's
// defaults for unset fields.
//
// The values are marshalled rather than concatenated: a description containing
// a colon - or starting with a "#" - would otherwise produce invalid YAML.
func writeSkillFrontmatter(b *strings.Builder, meta SkillMetadata) {
	if meta.Name == "" {
		meta.Name = defaultSkillName
	}

	if meta.Description == "" {
		meta.Description = defaultSkillDescription
	}

	b.WriteString("---\n")

	// Marshal cannot fail for a struct of strings, so a failure here leaves the
	// frontmatter empty rather than taking the whole skill down.
	if out, err := yaml.Marshal(meta); err == nil {
		b.Write(out)
	}

	b.WriteString("---\n\n")
}

// writeSkillDetails renders a command's arguments and flags as an indented list.
func writeSkillDetails(b *strings.Builder, c CommandInfo, pad string) {
	for _, a := range c.Arguments {
		name := util.Pick(a.Optional, "<"+a.Name+">", "["+a.Name+"]")

		b.WriteString(pad + "- `" + name + "`")

		if a.Repeat {
			b.WriteString(" (repeatable)")
		}

		if a.Description != "" {
			b.WriteString(" " + a.Description)
		}

		b.WriteString("\n")
	}

	for _, f := range c.Flags {
		fmt.Fprintf(b, "%s- `--%s` (%s) %s\n", pad, f.Name, f.Type, f.Description)
	}
}

// writeSkillSubcommands renders a command's descendants as an indented list. An
// agent needs the runnable leaf paths, which a flat listing cannot show.
func writeSkillSubcommands(b *strings.Builder, commands []CommandInfo, depth int) {
	for _, c := range commands {
		pad := strings.Repeat("  ", depth)

		b.WriteString(pad + "- `" + c.FullPath)

		if usage := c.Usage(); usage != "" {
			b.WriteString(" " + usage)
		}

		b.WriteString("`")

		if c.Description != "" {
			b.WriteString(" - " + c.Description)
		}

		b.WriteString("\n")

		writeSkillDetails(b, c, pad+"  ")

		writeSkillSubcommands(b, c.Subcommands, depth+1)
	}

	if depth == 0 && len(commands) > 0 {
		b.WriteString("\n")
	}
}

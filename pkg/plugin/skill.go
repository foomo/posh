package plugin

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// DefaultSkillsPath is the skills directory `posh agent skill install` writes
// into unless given an explicit path.
const DefaultSkillsPath = ".claude/skills"

// Skill directory naming. The root skill is RootSkillName; every command gets
// CommandSkillPrefix + its name.
//
// The prefix is the uninstall contract: `agent skill uninstall` removes what
// matches it and nothing else, so a hand-written skill sitting in the same
// directory survives. Keeping it a constant shared by writer and remover is what
// stops the two from drifting - a manifest would be the alternative, but it is
// extra on-disk state that can itself go stale.
const (
	RootSkillName      = "posh"
	CommandSkillPrefix = "posh-"
)

// Defaults for the generated frontmatter, used for any SkillMetadata field a
// project leaves unset.
//
// The description names triggering conditions rather than describing what posh
// is: it is the only thing an agent runtime sees when deciding whether to load
// the skill, so a summary of the tool ("drive this project's posh shell") gives
// it nothing to match a request against and the skill never fires.
const (
	defaultSkillName        = RootSkillName
	defaultSkillDescription = "Use when running any task in this project - building, testing, " +
		"deploying, or working with its clusters and services - and when asking how a task is " +
		"run here at all. This project drives such tasks through posh (Project Oriented Shell) " +
		"rather than ad-hoc shell commands."
)

// skillDetails points at the two places that carry a command's full argument and
// flag detail.
//
// The detail is deliberately not inlined. It is large - in a real project the
// flag and argument tree was over a quarter of the generated skill - and mostly
// redundant, because every leaf repeats the flag set it inherits. Both escape
// hatches below already render it on demand, and neither costs the agent
// anything until it asks.
const skillDetails = "Flags and arguments: `posh agent catalog` (JSON), or `posh help <command>`.\n\n"

// skillSetup documents posh's own subcommands, as opposed to the project
// commands the catalog is generated from.
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

// skillConventions states once what would otherwise be repeated in every
// command's own skill.
//
// Provider prose converged on the same handful of sentences - one real project
// repeated "that URL is also the schema's `$id`" twenty times across its
// fragments. Stating them here costs the agent one read of the root skill
// instead of one per command it loads.
const skillConventions = "## Conventions\n\n" +
	"These hold across this project's commands, so the individual command skills\n" +
	"do not repeat them:\n\n" +
	"- A schema URL in a config file doubles as that schema's `$id`.\n" +
	"- Defaults shown in a command's skill are overridable via `.posh.yaml`; that\n" +
	"  file is the authority when the two disagree.\n" +
	"- Confirm a value against both `.posh.yaml` and the command's own output\n" +
	"  before relying on it.\n\n"

// SkillCommand pairs a described command with the extra markdown and
// frontmatter its command.Skiller and command.SkillMetadataer implementations
// contributed, if any.
//
// The catalog shape stays untouched: skill prose is a rendering concern, so it
// rides alongside CommandInfo rather than inside it.
type SkillCommand struct {
	CommandInfo

	// Skill is the markdown contributed by command.Skiller, or empty.
	Skill string

	// Metadata is the frontmatter contributed by command.SkillMetadataer. Unset
	// fields fall back to values derived from the command.
	Metadata SkillMetadata
}

// SkillName is the directory name this command's skill is generated into.
func (c SkillCommand) SkillName() string {
	if c.Metadata.Name != "" {
		return c.Metadata.Name
	}

	// Nested paths ("cache clear") cannot appear here - only top level commands
	// get their own skill - but a name is still a path token, so anything
	// directory-unsafe is flattened rather than trusted.
	return CommandSkillPrefix + strings.ReplaceAll(c.FullPath, " ", "-")
}

// Describes reports whether a command says anything a skill of its own is worth
// spending tokens on: prose from command.Skiller, or a subcommand tree whose
// runnable leaf paths a flat index cannot show.
//
// Arguments and flags deliberately do not count. A skill never prints them - it
// points at `posh agent catalog` and `posh help` instead, see skillDetails - so
// a command whose only structure is an argument list would get a file saying
// nothing its index entry does not already say.
//
// The test is deliberately structural rather than "does it implement the
// interfaces". command.Describe already collapses a non-Describer into a
// CommandInfo carrying only a path and description, so by the time it arrives
// here a command that never opted in is indistinguishable from a Describer that
// legitimately describes a bare leaf - and both are equally uninformative.
//
// A command that fails the test still appears in the root skill's index with its
// name and description; it just does not get a skill of its own.
func (c SkillCommand) Describes() bool {
	return strings.TrimSpace(c.Skill) != "" || len(c.Subcommands) > 0
}

// RenderRootSkill renders the root SKILL.md: how to invoke posh, how to set it
// up, and an index of every command.
//
// It deliberately carries no per-command detail and no hazard prose. Nothing
// routes off it - each command that has more to say gets its own skill, which an
// agent loads on its own - so the root stays small enough to be worth loading
// unconditionally.
func RenderRootSkill(meta SkillMetadata, commands []SkillCommand) string {
	var b strings.Builder

	writeSkillFrontmatter(&b, meta, defaultSkillName, defaultSkillDescription)

	b.WriteString("# posh\n\n")
	b.WriteString("Run project tasks via `posh execute <command...>` (alias `posh x`) instead of\n")
	b.WriteString("guessing shell commands directly. JSON output is automatic under an agent\n")
	b.WriteString("harness; `POSH_AGENT_MODE=0` forces human output. Re-run\n")
	b.WriteString("`posh agent skill update` after this project's commands change to refresh\n")
	b.WriteString("these files.\n\n")
	b.WriteString("posh does not enforce access control: whatever the surrounding harness or CI\n")
	b.WriteString("permits is the real boundary. A command being listed here does not mean it is\n")
	b.WriteString("safe to run unattended.\n\n")
	b.WriteString(skillSetup)
	b.WriteString(skillConventions)
	b.WriteString("## Commands\n\n")
	b.WriteString("Run `posh agent catalog` for every command as JSON, including arguments and\n")
	b.WriteString("flags, or `posh help` for the human-readable list.\n\n")

	for _, c := range commands {
		fmt.Fprintf(&b, "- `%s`", c.FullPath)

		if c.Description != "" {
			b.WriteString(" - " + c.Description)
		}

		// Commands with more to say have a skill of their own; name it, so an
		// agent reading the index knows there is more and where it is.
		if c.Describes() {
			fmt.Fprintf(&b, " (skill: `%s`)", c.SkillName())
		}

		b.WriteString("\n")
	}

	return b.String()
}

// RenderCommandSkill renders one command's own SKILL.md: its runnable paths and
// whatever prose it contributes.
func RenderCommandSkill(c SkillCommand) string {
	var b strings.Builder

	writeSkillFrontmatter(&b, c.Metadata, c.SkillName(), defaultCommandSkillDescription(c))

	heading := c.FullPath
	if usage := c.Usage(); usage != "" {
		heading += " " + usage
	}

	fmt.Fprintf(&b, "# `%s`\n\n", heading)

	if c.Description != "" {
		b.WriteString(c.Description + "\n\n")
	}

	b.WriteString("Run it as `posh execute " + c.FullPath + "` (alias `posh x`). " +
		"See the `posh` skill\nfor setup and project-wide conventions.\n\n")
	b.WriteString(skillDetails)

	writeSkillSubcommands(&b, c.Subcommands, 0)

	// Trimmed so a contribution ending in a newline does not stack blank lines
	// on top of the one added here.
	if skill := strings.TrimSpace(c.Skill); skill != "" {
		b.WriteString(skill + "\n")
	}

	return b.String()
}

// defaultCommandSkillDescription derives a description for a command that did
// not supply one.
//
// It is a fallback, not a good description: it can only restate the one-line
// description, which says what the command is rather than when to reach for it.
// `agent skill install` reports every command that lands here, see
// FallbackSkillDescriptions.
func defaultCommandSkillDescription(c SkillCommand) string {
	ret := "Use when running `" + c.FullPath + "` commands in this project."

	if c.Description != "" {
		ret += " " + strings.ToUpper(c.Description[:1]) + c.Description[1:]

		if !strings.HasSuffix(ret, ".") {
			ret += "."
		}
	}

	return ret
}

// FallbackSkillDescriptions returns the paths of the commands that did not
// supply a frontmatter description, so a caller can report the gap rather than
// let it pass silently.
//
// The description is the only thing an agent runtime matches a request against,
// so a derived one usually means that command's skill never loads.
func FallbackSkillDescriptions(commands []SkillCommand) []string {
	var ret []string

	for _, c := range commands {
		if c.Describes() && c.Metadata.Description == "" {
			ret = append(ret, c.FullPath)
		}
	}

	return ret
}

// writeSkillFrontmatter renders the YAML frontmatter, falling back to the given
// defaults for unset fields.
//
// The values are marshalled rather than concatenated: a description containing
// a colon - or starting with a "#" - would otherwise produce invalid YAML.
func writeSkillFrontmatter(b *strings.Builder, meta SkillMetadata, name, description string) {
	if meta.Name == "" {
		meta.Name = name
	}

	if meta.Description == "" {
		meta.Description = description
	}

	b.WriteString("---\n")

	// Marshal cannot fail for a struct of strings, so a failure here leaves the
	// frontmatter empty rather than taking the whole skill down.
	if out, err := yaml.Marshal(meta); err == nil {
		b.Write(out)
	}

	b.WriteString("---\n\n")
}

// writeSkillSubcommands renders a command's descendants as an indented list. An
// agent needs the runnable leaf paths, which a flat listing cannot show.
//
// Arguments and flags are left out on purpose, see skillDetails.
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

		writeSkillSubcommands(b, c.Subcommands, depth+1)
	}

	if depth == 0 && len(commands) > 0 {
		b.WriteString("\n")
	}
}

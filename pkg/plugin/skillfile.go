package plugin

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// List asks plg for its command catalog, requiring the optional Lister
// interface. It returns a descriptive error when the plugin does not implement
// it, rather than leaving callers to repeat the type assertion.
//
// Resolving the plugin itself is the caller's job: the posh CLI holds its
// provider in package state, and a downstream shell may have one in hand
// already.
func List(ctx context.Context, plg any) ([]CommandInfo, error) {
	lister, ok := plg.(Lister)
	if !ok {
		return nil, errors.New("this project's posh shell does not support the agent command catalog")
	}

	return lister.List(ctx), nil
}

// ListSkill asks plg for its command catalog together with the extra SKILL.md
// markdown each command contributes.
//
// A plugin implementing only Lister falls back to a catalog without prose, so
// the skill still renders for shells that have not opted in.
func ListSkill(ctx context.Context, plg any) ([]SkillCommand, error) {
	if lister, ok := plg.(SkillLister); ok {
		return lister.ListSkill(ctx), nil
	}

	commands, err := List(ctx, plg)
	if err != nil {
		return nil, err
	}

	ret := make([]SkillCommand, 0, len(commands))
	for _, value := range commands {
		ret = append(ret, SkillCommand{CommandInfo: value})
	}

	return ret, nil
}

// SkillMetadataOf returns the root SKILL.md frontmatter plg supplies, or the
// zero value - which renders posh's defaults - when it does not implement
// SkillMetadataer.
//
// This is the plugin-level frontmatter, covering the root skill only. A command
// supplies its own through command.SkillMetadataer.
func SkillMetadataOf(ctx context.Context, plg any) SkillMetadata {
	if value, ok := plg.(SkillMetadataer); ok {
		return value.SkillMetadata(ctx)
	}

	return SkillMetadata{}
}

// WriteSkill generates the root skill and one skill per command into the skills
// directory at path, creating it if needed. Passing an empty path writes to
// DefaultSkillsPath.
//
// One skill per command rather than one file for all of them: a single file
// means an agent pays for every command in the project in order to use one. Only
// commands that Describes() get their own; the rest are covered by the root
// skill's index.
//
// Previously generated skills are removed first, so a command that was renamed
// or dropped since the last run leaves nothing behind - regenerating is what
// keeps the result honest, and a stale directory would otherwise survive
// indefinitely because nothing ever revisits it. Hand-written skills alongside
// them are untouched, see CommandSkillPrefix.
//
// It returns the paths written, in the order written.
func WriteSkill(path string, meta SkillMetadata, commands []SkillCommand) ([]string, error) {
	root := SkillsPath(path)

	if err := RemoveSkill(root); err != nil {
		return nil, err
	}

	ret := []string{filepath.Join(root, RootSkillName, "SKILL.md")}

	contents := []string{RenderRootSkill(meta, commands)}

	for _, c := range commands {
		if !c.Describes() {
			continue
		}

		ret = append(ret, filepath.Join(root, c.SkillName(), "SKILL.md"))
		contents = append(contents, RenderCommandSkill(c))
	}

	for i, name := range ret {
		if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
			return nil, errors.Wrapf(err, "failed to create directory for %s", name)
		}

		if err := os.WriteFile(name, []byte(contents[i]), 0o644); err != nil { //nolint:gosec
			return nil, errors.Wrapf(err, "failed to write %s", name)
		}
	}

	return ret, nil
}

// RemoveSkill deletes every generated skill directory under the skills
// directory at path, falling back to DefaultSkillsPath when empty.
//
// Only the root skill and directories carrying CommandSkillPrefix are removed,
// so a hand-written skill in the same directory survives. A missing directory is
// not an error, so uninstall is idempotent - and it deliberately needs no
// config or resolved plugin, so cleanup still works when the shell itself is
// broken.
func RemoveSkill(path string) error {
	root := SkillsPath(path)

	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return errors.Wrapf(err, "failed to read %s", root)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if name := entry.Name(); name != RootSkillName && !strings.HasPrefix(name, CommandSkillPrefix) {
			continue
		}

		if err := os.RemoveAll(filepath.Join(root, entry.Name())); err != nil {
			return errors.Wrapf(err, "failed to remove %s", filepath.Join(root, entry.Name()))
		}
	}

	return nil
}

// SkillsPath resolves a user supplied skills directory, falling back to
// DefaultSkillsPath when empty.
func SkillsPath(path string) string {
	if path == "" {
		return DefaultSkillsPath
	}

	return path
}

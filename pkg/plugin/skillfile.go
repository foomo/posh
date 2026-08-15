package plugin

import (
	"context"
	"os"
	"path/filepath"

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

// SkillMetadataOf returns the SKILL.md frontmatter plg supplies, or the zero
// value - which renders posh's defaults - when it does not implement
// SkillMetadataer.
func SkillMetadataOf(ctx context.Context, plg any) SkillMetadata {
	if value, ok := plg.(SkillMetadataer); ok {
		return value.SkillMetadata(ctx)
	}

	return SkillMetadata{}
}

// WriteSkill renders commands as a Claude Code SKILL.md and writes it to path,
// creating any missing parent directories.
//
// Passing an empty path writes to DefaultSkillPath. The file is regenerated in
// full every time, so re-running after a project's commands change always
// produces an up-to-date skill.
func WriteSkill(path string, meta SkillMetadata, commands []SkillCommand) error {
	path = SkillPath(path)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return errors.Wrapf(err, "failed to create directory for %s", path)
	}

	return os.WriteFile(path, []byte(RenderSkill(meta, commands)), 0o644) //nolint:gosec
}

// RemoveSkill deletes the skill file at path, falling back to DefaultSkillPath
// when empty. A missing file is not an error, so uninstall is idempotent.
func RemoveSkill(path string) error {
	path = SkillPath(path)

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "failed to remove %s", path)
	}

	return nil
}

// SkillPath resolves a user supplied skill path, falling back to
// DefaultSkillPath when empty.
func SkillPath(path string) string {
	if path == "" {
		return DefaultSkillPath
	}

	return path
}

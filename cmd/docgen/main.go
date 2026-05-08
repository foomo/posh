// Command docgen generates Markdown documentation for every cobra command
// shipped by posh — both the binary-only commands (init, config, version) and
// the downstream-only commands that require a Plugin (prompt, execute, brew,
// require).
//
// Output is one Markdown file per command, suitable for VitePress. The file
// header is rewritten so each page has a `# posh <subcommand>` title and a
// short YAML front matter block. Internal links between parent and subcommand
// pages keep cobra's default `.md` extension — VitePress resolves them.
//
// Usage:
//
//	go run ./cmd/docgen --output docs/reference/cli
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/foomo/posh/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	out := flag.String("output", "docs/reference/cli", "output directory")
	flag.Parse()

	if err := os.MkdirAll(*out, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	root := buildFullTree()

	prepender := func(filename string) string {
		base := strings.TrimSuffix(filepath.Base(filename), ".md")
		title := strings.ReplaceAll(base, "_", " ")
		return fmt.Sprintf("---\ntitle: %s\n---\n\n", title)
	}

	linkHandler := func(name string) string { return name }

	if err := doc.GenMarkdownTreeCustom(root, *out, prepender, linkHandler); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if err := escapeAngleBrackets(*out); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Println("wrote markdown to", *out)
}

// escapeAngleBrackets walks every Markdown file under dir and HTML-escapes
// `<` and `>` characters that appear outside fenced code blocks. VitePress
// renders Markdown through Vue's template parser, which trips on bare
// placeholders like `<MODNAME>` lifted from cobra's Long descriptions.
func escapeAngleBrackets(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		if !strings.HasPrefix(e.Name(), "posh") {
			continue
		}

		path := filepath.Join(dir, e.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var b strings.Builder
		inFence := false
		for _, line := range strings.SplitAfter(string(raw), "\n") {
			trimmed := strings.TrimLeft(line, " \t")
			if strings.HasPrefix(trimmed, "```") {
				inFence = !inFence
				b.WriteString(line)
				continue
			}

			if inFence {
				b.WriteString(line)
				continue
			}

			line = strings.ReplaceAll(line, "<", "&lt;")
			line = strings.ReplaceAll(line, ">", "&gt;")
			b.WriteString(line)
		}

		if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
			return err
		}
	}

	return nil
}

// buildFullTree assembles a fresh cobra root with **all** subcommands, both
// the standalone-binary ones and the downstream-only ones. It bypasses
// cmd.Init's plugin-conditional logic on purpose: cobra/doc only inspects
// Use/Short/Long/flags, so RunE never fires and the nil pluginProvider is
// safe.
func buildFullTree() *cobra.Command {
	root := cmd.NewRoot()
	cmd.NewInit(root)
	cmd.NewConfig(root)
	cmd.NewVersion(root)
	cmd.NewBrew(root)
	cmd.NewExecute(root)
	cmd.NewPrompt(root)
	cmd.NewRequire(root)

	root.DisableAutoGenTag = true
	return root
}

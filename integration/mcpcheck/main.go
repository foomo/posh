// Command mcpcheck drives a built posh binary's `mcp` subcommand over stdio,
// calling list_commands and run_command, and exits non-zero if either call
// fails or returns an unexpected shape. It exists for make test.demo, which
// otherwise only exercises `posh execute` directly - not the MCP server this
// binary also serves once a provider is wired up.
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: mcpcheck <path-to-posh-binary>")
		os.Exit(1)
	}

	if err := run(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(binary string) error {
	ctx := context.Background()

	binary, err := filepath.Abs(binary)
	if err != nil {
		return fmt.Errorf("resolving binary path: %w", err)
	}

	// posh mcp's PreRunE loads .posh.yaml from the current working directory,
	// same as execute/prompt/require - the binary lives at <project>/bin/posh,
	// so the project root is two directories up.
	projectDir := filepath.Dir(filepath.Dir(binary))

	client := mcp.NewClient(&mcp.Implementation{Name: "mcpcheck", Version: "1"}, nil)
	cmd := exec.CommandContext(ctx, binary, "mcp") //nolint:gosec // binary is a trusted local path from make test.demo's own build output, not external input
	cmd.Dir = projectDir
	transport := &mcp.CommandTransport{Command: cmd}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer session.Close()

	listRes, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "list_commands"})
	if err != nil {
		return fmt.Errorf("list_commands: %w", err)
	}

	if listRes.IsError {
		return fmt.Errorf("list_commands returned an error result: %+v", listRes.Content)
	}

	runRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "run_command",
		Arguments: map[string]any{"args": []string{"welcome", "demo"}},
	})
	if err != nil {
		return fmt.Errorf("run_command: %w", err)
	}

	if runRes.IsError {
		return fmt.Errorf("run_command returned an error result: %+v", runRes.Content)
	}

	fmt.Println("mcpcheck: list_commands and run_command both succeeded")

	return nil
}

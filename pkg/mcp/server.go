package mcp

import (
	"context"
	"encoding/json"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/plugin"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server serves this project's posh shell as an MCP server over stdio.
type Server struct {
	l        log.Logger
	provider plugin.Provider
}

// New constructs a Server. provider resolves a fresh plugin instance per
// call, each wired to that call's own logger - required so a later
// run_command tool (see run_command.go) can give each invocation a private
// buffering logger rather than sharing one resolved plugin.Plugin across
// calls.
func New(l log.Logger, provider plugin.Provider) *Server {
	return &Server{l: l, provider: provider}
}

// Run serves the MCP protocol over stdio until ctx is canceled or the client
// disconnects.
func (s *Server) Run(ctx context.Context) error {
	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "posh", Version: "1"}, nil)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "posh_list_commands",
		Description: "List this project's posh commands, the same catalog `posh agent catalog` prints.",
	}, s.listCommands)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "posh_run_command",
		Description: "Run a posh command, the same as `posh execute <args...>`, and return its output.",
	}, s.runCommand)

	return server.Run(ctx, &sdkmcp.StdioTransport{})
}

type listCommandsInput struct{}

// listCommandsOutput is left empty deliberately: plugin.CommandInfo nests
// itself through Subcommands, and the SDK's AddTool derives a JSON schema for
// Out by reflection, which panics on that cycle. The catalog is emitted as
// unstructured TextContent (JSON) instead of a typed Out.
type listCommandsOutput struct{}

func (s *Server) listCommands(ctx context.Context, req *sdkmcp.CallToolRequest, input listCommandsInput) (*sdkmcp.CallToolResult, listCommandsOutput, error) {
	plg, err := s.provider(s.l)
	if err != nil {
		return errorResult(err), listCommandsOutput{}, nil
	}

	catalog, err := ListCommands(ctx, plg)
	if err != nil {
		return errorResult(err), listCommandsOutput{}, nil
	}

	out, err := json.Marshal(catalog)
	if err != nil {
		return errorResult(err), listCommandsOutput{}, nil
	}

	return &sdkmcp.CallToolResult{Content: []sdkmcp.Content{
		&sdkmcp.TextContent{Text: string(out)},
	}}, listCommandsOutput{}, nil
}

func errorResult(err error) *sdkmcp.CallToolResult {
	return &sdkmcp.CallToolResult{IsError: true, Content: []sdkmcp.Content{
		&sdkmcp.TextContent{Text: err.Error()},
	}}
}

type runCommandInput struct {
	Args []string `json:"args" jsonschema:"the posh command and its arguments, e.g. [\"cache\", \"clear\"]"`
}

type runCommandOutput struct {
	Output string `json:"output"`
}

func (s *Server) runCommand(ctx context.Context, req *sdkmcp.CallToolRequest, input runCommandInput) (*sdkmcp.CallToolResult, runCommandOutput, error) {
	output, err := RunCommand(ctx, s.provider, input.Args)
	if err != nil {
		result := errorResult(err)
		result.Content = append(result.Content, &sdkmcp.TextContent{Text: output})

		return result, runCommandOutput{}, nil
	}

	return nil, runCommandOutput{Output: output}, nil
}

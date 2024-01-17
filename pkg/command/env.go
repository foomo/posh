package command

import (
	"context"
	"os"
	"sort"
	"strings"

	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/pterm/pterm"
)

type Env struct {
	l    log.Logger
	tree tree.Root
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewEnv(l log.Logger) *Env {
	inst := &Env{
		l: l,
	}
	inst.tree = tree.New(&tree.Node{
		Name:        "env",
		Description: "Manage internal environment variables",
		Nodes: tree.Nodes{
			{
				Name:        "list",
				Description: "List all environment variables",
				Execute:     inst.list,
			},
			{
				Name:        "set",
				Description: "Set an internal environment variable",
				Args: tree.Args{
					{
						Name:        "Key",
						Description: "Key of the environment variable.",
					},
					{
						Name:        "Value",
						Optional:    true,
						Description: "Value of the environment variable.",
					},
				},
				Execute: inst.set,
			},
			{
				Name:        "unset",
				Description: "Unset an environment variable",
				Args: tree.Args{
					{
						Name:        "Key",
						Description: "Key of the environment variable.",
					},
				},
				Execute: inst.unset,
			},
		},
	})
	return inst
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (c *Env) Name() string {
	return c.tree.Node().Name
}

func (c *Env) Description() string {
	return c.tree.Node().Description
}

func (c *Env) Complete(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
	return c.tree.Complete(ctx, r)
}

func (c *Env) Execute(ctx context.Context, r *readline.Readline) error {
	return c.tree.Execute(ctx, r)
}

func (c *Env) Help(ctx context.Context, r *readline.Readline) string {
	return c.tree.Help(ctx, r)
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (c *Env) set(ctx context.Context, r *readline.Readline) error {
	return os.Setenv(r.Args().At(1), r.Args().AtDefault(2, ""))
}

func (c *Env) unset(ctx context.Context, r *readline.Readline) error {
	return os.Unsetenv(r.Args().At(1))
}

func (c *Env) list(ctx context.Context, r *readline.Readline) error {
	data := pterm.TableData{{"Name", "Value"}}
	values := os.Environ()
	sort.Strings(values)
	for _, s := range values {
		data = append(data, strings.SplitN(s, "=", 2))
	}
	return pterm.DefaultTable.WithHasHeader(true).WithData(data).Render()
}

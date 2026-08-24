package command

import (
	"context"
	"os"
	"sort"
	"strings"

	"github.com/foomo/posh/pkg/agent"
	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/foomo/posh/pkg/util/suggests"
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
				Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
					fs.Default().Bool("with-values", false, "include values in the list")
					return nil
				},
				Execute: inst.list,
			},
			{
				Name:        "get",
				Description: "Get an environment variable",
				Args: tree.Args{
					{
						Name:        "key",
						Description: "Name of the environment variable",
						Suggest: func(ctx context.Context, t tree.Root, r *readline.Readline) []goprompt.Suggest {
							return suggests.List(inst.envKeys())
						},
					},
				},
				Execute: inst.get,
			},
			{
				Name:        "set",
				Description: "Set an internal environment variable",
				Args: tree.Args{
					{
						Name:        "key",
						Description: "Name of the environment variable",
					},
					{
						Name:        "value",
						Optional:    true,
						Description: "Value of the environment variable",
					},
				},
				Execute: inst.set,
			},
			{
				Name:        "unset",
				Description: "Unset an environment variable",
				Args: tree.Args{
					{
						Name:        "key",
						Description: "Name of the environment variable",
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

// Describe implements the optional Describer interface, letting
// `posh agent catalog` describe this command's subtree.
func (c *Env) Describe(ctx context.Context) CommandInfo {
	return c.tree.Describe(ctx)
}

// Skill implements the optional Skiller interface. The scope of a `set` - this
// process only - is the thing an agent gets wrong, and no amount of structure
// conveys it.
func (c *Env) Skill(ctx context.Context, name string) string {
	return "#### Notes\n\n" +
		"These are the environment variables of the running shell process. `env set`\n" +
		"and `env unset` affect that process and the commands it goes on to run;\n" +
		"they do not touch your shell, `.env` files, or anything on disk, and the\n" +
		"change is gone when the shell exits.\n\n" +
		"That scope matters under `posh execute`: each invocation is a fresh\n" +
		"process, so `posh x env set FOO bar` followed by a separate\n" +
		"`posh x env get FOO` will not see the value. To set a variable for a\n" +
		"command, use the environment of the `posh` invocation itself\n" +
		"(`FOO=bar posh x ...`). Inside an interactive shell session the value does\n" +
		"persist across commands.\n\n" +
		"`env list` prints names only; pass `--with-values` to include values. Do\n" +
		"that deliberately - the environment routinely holds tokens and\n" +
		"credentials, and the output goes wherever the transcript goes."
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

func (c *Env) get(ctx context.Context, r *readline.Readline) error {
	name := r.Args().At(1)
	value := os.Getenv(name)

	return agent.Render(
		func() any { return EnvGet{Name: name, Value: value} },
		func() error {
			c.l.Info(value)
			return nil
		},
	)
}

func (c *Env) set(ctx context.Context, r *readline.Readline) error {
	return os.Setenv(r.Args().At(1), r.Args().AtDefault(2, ""))
}

func (c *Env) unset(ctx context.Context, r *readline.Readline) error {
	return os.Unsetenv(r.Args().At(1))
}

func (c *Env) list(ctx context.Context, r *readline.Readline) error {
	withValues, err := r.FlagSets().Default().GetBool("with-values")
	if err != nil {
		return err
	}

	values := os.Environ()
	sort.Strings(values)

	var pairs [][]string
	for _, s := range values {
		pairs = append(pairs, strings.SplitN(s, "=", 2))
	}

	return agent.Render(
		func() any { return c.listPayload(pairs, withValues) },
		func() error { return c.listTable(pairs, withValues) },
	)
}

// listPayload builds the JSON form. Unlike listTable it never wraps values to
// the terminal width, which would corrupt long entries such as PATH with
// embedded newlines.
func (c *Env) listPayload(pairs [][]string, withValues bool) EnvList {
	values := make([]EnvVar, 0, len(pairs))

	for _, pair := range pairs {
		value := EnvVar{Name: pair[0]}
		if withValues {
			value.Value = pair[1]
		}

		values = append(values, value)
	}

	return EnvList{Values: values}
}

// listTable renders the human-formatted table, wrapping values to the terminal
// width so long entries stay inside their column.
func (c *Env) listTable(pairs [][]string, withValues bool) error {
	if !withValues {
		data := pterm.TableData{{"Name"}}
		for _, pair := range pairs {
			data = append(data, []string{pair[0]})
		}

		return agent.Table(data, agent.WithHeader())
	}

	var maxKeyLen int
	for _, pair := range pairs {
		maxKeyLen = max(maxKeyLen, len(pair[0]))
	}

	maxValueLen := pterm.GetTerminalWidth() - maxKeyLen - 5

	for i, pair := range pairs {
		var value strings.Builder
		for len(pair[1]) > maxValueLen {
			value.WriteString(pair[1][:maxValueLen] + "\n")
			pair[1] = pair[1][maxValueLen:]
		}

		pairs[i][1] = value.String() + pair[1]
	}

	return agent.Table(append(pterm.TableData{{"Name", "Value"}}, pairs...), agent.WithHeader())
}

func (c *Env) envKeys() []string {
	var ret []string
	for _, s := range os.Environ() {
		ret = append(ret, strings.SplitN(s, "=", 2)[0])
	}

	return ret
}

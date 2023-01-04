package command

import (
	"context"
	"io/fs"
	"path"

	"github.com/charlievieth/fastwalk"
	"github.com/foomo/posh/pkg/cache"
	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/foomo/posh/pkg/shell"
	"github.com/foomo/posh/pkg/util/suggests"
)

// GoMod command
type GoMod struct {
	l      log.Logger
	cache  cache.Namespace
	parser *tree.Root
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

// NewGoMod command
func NewGoMod(l log.Logger, cache cache.Cache) *GoMod {
	inst := &GoMod{
		l:     l,
		cache: cache.Get("gomod"),
	}

	pathArg := &tree.Arg{
		Name:     "path",
		Optional: true,
		Suggest:  inst.completePaths,
	}

	inst.parser = &tree.Root{
		Name: "gomod",
		Nodes: []*tree.Node{
			{
				Name:        "tidy",
				Description: "run go mod tidy",
				Args:        []*tree.Arg{pathArg},
				Execute:     inst.tidy,
			},
			{
				Name:        "download",
				Description: "run go mod download",
				Args:        []*tree.Arg{pathArg},
				Execute:     inst.download,
			},
			{
				Name:        "outdated",
				Description: "show go mod outdated",
				Args:        []*tree.Arg{pathArg},
				Execute:     inst.outdated,
			},
		},
	}
	return inst
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (c *GoMod) Name() string {
	return c.parser.Name
}

func (c *GoMod) Description() string {
	return "run go mod"
}

func (c *GoMod) Complete(ctx context.Context, r *readline.Readline, d prompt.Document) []prompt.Suggest {
	return c.parser.Complete(ctx, r)
}

func (c *GoMod) Execute(ctx context.Context, r *readline.Readline) error {
	return c.parser.Execute(ctx, r)
}

func (c *GoMod) Help() string {
	return `Looks for go.mod files and runs the given command.

Usage:
  gomod [command] <path>

Available commands:
  tidy       run go mod tidy on specific or all paths
  download   run go mod download on specific or all paths
  outdated   list outdated packages on specific or all paths

Examples:
  gomod tidy ./path
`
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (c *GoMod) tidy(ctx context.Context, r *readline.Readline) error {
	var paths []string
	if r.Args().HasIndex(1) {
		paths = []string{r.Args().At(1)}
	} else {
		paths = c.paths()
	}
	for _, value := range paths {
		c.l.Info("go mod tidy:", value)
		if err := shell.New(ctx, c.l,
			"go", "mod", "tidy",
		).
			Args(r.AdditionalArgs()...).
			Dir(value).
			Run(); err != nil {
			return err
		}
	}
	return nil
}

func (c *GoMod) download(ctx context.Context, r *readline.Readline) error {
	var paths []string
	if r.Args().HasIndex(1) {
		paths = []string{r.Args().At(1)}
	} else {
		paths = c.paths()
	}
	for _, value := range paths {
		c.l.Info("go mod download:", value)
		if err := shell.New(ctx, c.l,
			"go", "mod", "tidy",
		).
			Args(r.AdditionalArgs()...).
			Dir(value).
			Run(); err != nil {
			return err
		}
	}
	return nil
}

func (c *GoMod) outdated(ctx context.Context, r *readline.Readline) error {
	var paths []string
	if r.Args().HasIndex(1) {
		paths = []string{r.Args().At(1)}
	} else {
		paths = c.paths()
	}
	for _, value := range paths {
		c.l.Info("go mod outdated:", value)
		if err := shell.New(ctx, c.l,
			"go", "list",
			"-u", "-m", "-json", "all",
			"|", "go-mod-outdated", "-update", "-direct",
		).
			Args(r.AdditionalArgs()...).
			Dir(value).
			Run(); err != nil {
			return err
		}
	}
	return nil
}

func (c *GoMod) completePaths(ctx context.Context, p *tree.Root, r *readline.Readline) []prompt.Suggest {
	return suggests.List(c.paths())
}

//nolint:forcetypeassert
func (c *GoMod) paths() []string {
	return c.cache.Get("paths", func() any {
		var ret []string
		if err := fastwalk.Walk(&fastwalk.Config{
			Follow: false,
		}, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.Name() == "go.mod" {
				ret = append(ret, path.Dir(p))
			}
			return nil
		}); err != nil {
			c.l.Debug("failed to walk files", err.Error())
		}
		return ret
	}).([]string)
}

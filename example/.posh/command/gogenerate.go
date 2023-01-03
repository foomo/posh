package command

import (
	"context"
	"io/fs"
	"os"

	"github.com/charlievieth/fastwalk"
	"github.com/foomo/posh/pkg/cache"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/foomo/posh/pkg/shell"
	"github.com/pkg/errors"
)

// GoGenerate command
type GoGenerate struct {
	l     log.Logger
	name  string
	cache cache.Namespace
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

// NewGoGenerate command
func NewGoGenerate(l log.Logger, cache cache.Cache) *GoGenerate {
	return &GoGenerate{
		l:     l,
		name:  "gogenerate",
		cache: cache.Get("gogenerate"),
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (c *GoGenerate) Name() string {
	return c.name
}

func (c *GoGenerate) Description() string {
	return "run go generate on generate.go files"
}

func (c *GoGenerate) Complete(ctx context.Context, r *readline.Readline, d prompt.Document) (suggests []prompt.Suggest) {
	if r.Args().LenLte(1) {
		suggests = c.completePaths()
	}
	return nil
}

func (c *GoGenerate) Validate(ctx context.Context, r *readline.Readline) error {
	switch {
	case r.Args().LenIs(0):
		return nil
	case r.Args().LenGt(1):
		return errors.New("too many parameters")
	}
	if info, err := os.Stat(r.Args().At(1)); err != nil || info.IsDir() {
		return errors.Errorf("invalid [path] parameter: %s", r.Args().At(1))
	}
	return nil
}

func (c *GoGenerate) Execute(ctx context.Context, r *readline.Readline) error {
	var paths []string
	if r.Args().HasIndex(0) {
		paths = append(paths, r.Args().At(0))
	} else {
		paths = c.paths()
	}

	for _, value := range paths {
		c.l.Info("go generate:", value)
		if err := shell.New(ctx, c.l,
			"go", "generate", value,
		).
			Args(r.AdditionalArgs()...).
			Run(); err != nil {
			return err
		}
	}
	return nil
}

func (c *GoGenerate) Help() string {
	return `Looks for generate.go files and runs them.

Usage:
  gogenerate [path]

Examples:
  gogenerate tidy ./path/to/generate.go
`
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (c *GoGenerate) completePaths() (suggest []prompt.Suggest) {
	for _, value := range c.paths() {
		suggest = append(suggest, prompt.Suggest{Text: value})
	}
	return suggest
}

func (c *GoGenerate) paths() []string {
	return c.cache.Get("paths", func() interface{} {
		var ret []string
		if err := fastwalk.Walk(&fastwalk.Config{
			Follow: false,
		}, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.Name() == "generate.go" {
				ret = append(ret, path)
			}
			return nil
		}); err != nil {
			c.l.Debug("failed to walk files", err.Error())
		}
		return ret
	}).([]string)
}

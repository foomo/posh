package scaffold

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/foomo/posh/pkg/log"
	"github.com/pkg/errors"
)

type (
	Scaffold struct {
		l     log.Logger
		dry   bool
		force bool
	}
	Option func(*Scaffold) error
)

// ------------------------------------------------------------------------------------------------
// ~ Options
// ------------------------------------------------------------------------------------------------

func WithDry(v bool) Option {
	return func(o *Scaffold) error {
		o.dry = v
		return nil
	}
}

func WithForce(v bool) Option {
	return func(o *Scaffold) error {
		o.force = v
		return nil
	}
}

func WithLogger(v log.Logger) Option {
	return func(o *Scaffold) error {
		o.l = v
		return nil
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func New(opts ...Option) (*Scaffold, error) {
	inst := &Scaffold{
		l:     log.NewFmt(),
		dry:   false,
		force: false,
	}
	for _, opt := range opts {
		if opt != nil {
			if err := opt(inst); err != nil {
				return nil, err
			}
		}
	}
	return inst, nil
}

func (s *Scaffold) Render(source fs.FS, target string, vars any) error {
	// validate target
	if stat, err := os.Stat(target); errors.Is(err, os.ErrNotExist) {
		s.l.Print("scaffold:", target)
		if err := os.MkdirAll(target, os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to create target folder (%s)", target)
		}
	} else if err != nil {
		return errors.Wrapf(err, "failed to stat target folder (%s)", target)
	} else if !stat.IsDir() {
		return fmt.Errorf("target not a directory (%s)", target)
	}

	// iterate source
	if err := fs.WalkDir(source, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to walk fs dir")
		}

		filename := filepath.Join(target, strings.ReplaceAll(path, "$", ""))

		if path == "." {
			return nil
		} else if d.IsDir() {
			s.l.Print("scaffold:", filename)
			return os.MkdirAll(filename, os.ModePerm)
		} else {
			tpl, err := template.New(d.Name()).Funcs(sprig.FuncMap()).ParseFS(source, path)
			if err != nil {
				return err
			}

			if stat, err := os.Stat(filename); errors.Is(err, fs.ErrNotExist) {
				// all good
			} else if err != nil {
				return errors.Wrapf(err, "failed to stat target (%s)", filename)
			} else if stat.IsDir() {
				return fmt.Errorf("target file is an existing directory (%s)", filename)
			} else if !s.force {
				return fmt.Errorf("target file already exists (%s)", filename)
			}

			var out io.Writer
			if s.dry {
				out = os.Stdout
			} else {
				s.l.Print("scaffold:", filename)
				if file, err := os.Create(filename); err != nil {
					return errors.Wrapf(err, "failed to create target file (%s)", filename)
				} else {
					out = file
					defer func() {
						_ = file.Close()
					}()
				}
			}

			if err := tpl.Execute(out, vars); err != nil {
				return errors.Wrapf(err, "failed to render target file (%s)", filename)
			}

			return nil
		}
	}); err != nil {
		return errors.Wrapf(err, "failed to render scaffold to %s", target)
	}

	return nil
}

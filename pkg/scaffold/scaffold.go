package scaffold

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/alecthomas/chroma/quick"
	"github.com/foomo/posh/pkg/log"
	"github.com/pkg/errors"
)

type (
	Scaffold struct {
		l           log.Logger
		dry         bool
		override    bool
		directories []Directory
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

func WithOverride(v bool) Option {
	return func(o *Scaffold) error {
		o.override = v
		return nil
	}
}

func WithDirectories(v ...Directory) Option {
	return func(o *Scaffold) error {
		o.directories = append(o.directories, v...)
		return nil
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func New(l log.Logger, opts ...Option) (*Scaffold, error) {
	inst := &Scaffold{
		l:        l.Named("scaffold"),
		dry:      false,
		override: false,
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

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (s *Scaffold) Render(ctx context.Context) error {
	if err := s.renderDirectories(); err != nil {
		return err
	}
	return nil
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------
func (s *Scaffold) scaffoldDir(target string) error {
	s.l.Info("mkdir:", s.filename(target))
	if stat, err := os.Stat(target); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(target, os.ModePerm); err != nil {
			return err
		}
	} else if err != nil {
		return errors.Wrapf(err, "failed to stat target folder (%s)", target)
	} else if !stat.IsDir() {
		return fmt.Errorf("target not a directory (%s)", target)
	}
	return nil
}

func (s *Scaffold) scaffoldTemplate(target string, tpl *template.Template, data any) error {
	s.l.Info("file:", s.filename(target))
	file, err := os.Create(target)
	if err != nil {
		return errors.Wrapf(err, "failed to create target file (%s)", target)
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.l.Warn("failed to close file: %s", err.Error())
		}
	}()
	return tpl.Execute(file, data)
}

func (s *Scaffold) printTemplate(msg, target string, tpl *template.Template, data any) error {
	border := strings.Repeat("-", 80)
	s.l.Infof("\n%s\n%s: %s\n%s", border, msg, target, border)

	var out bytes.Buffer
	if err := tpl.Execute(&out, data); err != nil {
		return err
	}
	return quick.Highlight(os.Stdout, out.String(), filepath.Ext(target), "terminal", "monokai")
}

func (s *Scaffold) renderDirectories() error {
	for _, directory := range s.directories {
		if err := s.renderDirectory(directory); err != nil {
			return err
		}
	}
	return nil
}

func (s *Scaffold) renderDirectory(directory Directory) error {
	s.l.Info("scaffolding directory:", directory.Target)

	if err := s.scaffoldDir(directory.Target); err != nil {
		return err
	}

	if err := fs.WalkDir(directory.Source, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if path == "." {
			return nil
		} else if d.IsDir() {
			return s.scaffoldDir(s.filename(path))
		}

		filename := s.filename(path)

		tpl, err := template.New(d.Name()).Funcs(sprig.FuncMap()).ParseFS(directory.Source, path)
		if err != nil {
			return errors.Wrapf(err, "failed to parse source file (%s)", path)
		}

		if s.dry {
			return s.printTemplate("file", filename, tpl, directory.Data)
		} else if exists, err := s.fileExists(filename); err != nil {
			return s.printTemplate(err.Error(), filename, tpl, directory.Data)
		} else if exists && !s.override {
			return s.printTemplate("file exists", filename, tpl, directory.Data)
		} else {
			return s.scaffoldTemplate(filename, tpl, directory.Data)
		}
	}); err != nil {
		return errors.Wrapf(err, "failed to render scaffold to %s", directory.Target)
	}

	return nil
}

func (s *Scaffold) fileExists(target string) (bool, error) {
	if stat, err := os.Stat(target); errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, errors.Wrapf(err, "failed to stat target (%s)", target)
	} else if stat.IsDir() {
		return true, fmt.Errorf("target file is an existing directory (%s)", target)
	} else {
		return true, nil
	}
}

func (s *Scaffold) filename(v string) string {
	v = strings.ReplaceAll(v, "$", "")
	v = strings.TrimSuffix(v, ".gotext")
	v = strings.TrimSuffix(v, ".gohtml")
	return v
}

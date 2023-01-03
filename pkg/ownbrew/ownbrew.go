package ownbrew

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/log"
	"github.com/pkg/errors"
)

type (
	Ownbrew struct {
		l         log.Logger
		binDir    string
		tempDir   string
		caskDir   string
		cellarDir string
		packages  []config.Package
	}
	Option func(*Ownbrew) error
)

// ------------------------------------------------------------------------------------------------
// ~ Options
// ------------------------------------------------------------------------------------------------

func WithPackages(v ...config.Package) Option {
	return func(o *Ownbrew) error {
		o.packages = append(o.packages, v...)
		return nil
	}
}

func WithBinDir(v string) Option {
	return func(o *Ownbrew) error {
		o.binDir = v
		return nil
	}
}

func WithTempDir(v string) Option {
	return func(o *Ownbrew) error {
		o.tempDir = v
		return nil
	}
}

func WithCaskDir(v string) Option {
	return func(o *Ownbrew) error {
		o.caskDir = v
		return nil
	}
}

func WithCellarDir(v string) Option {
	return func(o *Ownbrew) error {
		o.cellarDir = v
		return nil
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func New(l log.Logger, opts ...Option) (*Ownbrew, error) {
	inst := &Ownbrew{
		l:         l,
		binDir:    "bin",
		tempDir:   ".posh/tmp",
		caskDir:   ".posh/scripts/ownbrew",
		cellarDir: ".posh/bin",
	}
	for _, opt := range opts {
		if opt != nil {
			if err := opt(inst); err != nil {
				return nil, err
			}
		}
	}

	for _, dir := range []string{inst.binDir, inst.tempDir, inst.caskDir, inst.cellarDir} {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}
	return inst, nil
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (o *Ownbrew) Install(ctx context.Context) error {
	o.l.Debug("installing packages", runtime.GOOS, runtime.GOARCH)

	for _, pkg := range o.packages {
		if pkg.Command == "" {
			pkg.Command = filepath.Join(o.caskDir, pkg.Name+".sh")
		}
		o.l.Debug("installing:", pkg.String())

		cellarFilename := filepath.Join(o.cellarDir, fmt.Sprintf("%s-%s-%s-%s", pkg.Name, pkg.Version, runtime.GOOS, runtime.GOARCH))
		exists, err := o.cellarExists(cellarFilename)
		if err != nil {
			return err
		}

		if !exists {
			if err := o.casksExists(pkg.Command); err != nil {
				return err
			}

			cmd := exec.CommandContext(ctx, "sh",
				pkg.Command,
				pkg.Version,
				runtime.GOOS,
				runtime.GOARCH,
			)

			cmd.Env = append(
				os.Environ(),
				"BIN_DIR="+o.cellarDir,
				"TEMP_DIR="+o.tempDir,
			)

			output, err := cmd.CombinedOutput()
			o.l.Debug(string(output))
			if err != nil {
				return errors.Wrapf(err, "failed to install package: %s", pkg.Name)
			}
		}

		// create symlink
		if err := o.symlink(cellarFilename, filepath.Join(o.binDir, pkg.Name)); err != nil {
			return err
		}
	}
	return nil
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (o *Ownbrew) symlink(source, target string) error {
	if err := os.Remove(target); os.IsNotExist(err) {
		// continue
	} else if err != nil {
		return err
	}

	prefix, err := filepath.Rel(filepath.Base(target), "")
	if err != nil {
		return err
	}
	prefix = strings.TrimSuffix(prefix, ".")

	o.l.Debug("symlink:", prefix+source, target)
	return os.Symlink(prefix+source, target)
}

func (o *Ownbrew) casksExists(filename string) error {
	if stat, err := os.Stat(filename); err != nil {
		return errors.Wrap(err, "failed to stat cask")
	} else if stat.IsDir() {
		return fmt.Errorf("not an executeable: %s", filename)
	} else {
		return nil
	}
}

func (o *Ownbrew) cellarExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		o.l.Debug("install:", filename)
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, "failed to stat bin")
	} else {
		o.l.Debug("exists:", filename)
		return true, nil
	}
}

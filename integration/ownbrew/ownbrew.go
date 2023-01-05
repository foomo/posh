package ownbrew

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/foomo/posh/pkg/log"
	"github.com/pkg/errors"
)

type (
	Ownbrew struct {
		l         log.Logger
		dry       bool
		binDir    string
		tapDir    string
		tempDir   string
		cellarDir string
		packages  []Package
	}
	Option func(*Ownbrew) error
)

// ------------------------------------------------------------------------------------------------
// ~ Options
// ------------------------------------------------------------------------------------------------

func WithDry(v bool) Option {
	return func(o *Ownbrew) error {
		o.dry = v
		return nil
	}
}

func WithPackages(v ...Package) Option {
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

func WithTapDir(v string) Option {
	return func(o *Ownbrew) error {
		o.tapDir = v
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
		l:         l.Named("ownbrew"),
		binDir:    "bin",
		tempDir:   ".posh/tmp",
		tapDir:    ".posh/ownbrew",
		cellarDir: ".posh/bin",
	}
	for _, opt := range opts {
		if opt != nil {
			if err := opt(inst); err != nil {
				return nil, err
			}
		}
	}

	for _, dir := range []string{inst.binDir, inst.tempDir, inst.tapDir, inst.cellarDir} {
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
	o.l.Debug("install:", runtime.GOOS, runtime.GOARCH)

	for _, pkg := range o.packages {
		cellarFilename := o.cellarFilename(pkg)
		if cellarExists, err := o.cellarExists(cellarFilename); err != nil {
			return err
		} else if !cellarExists {
			if pkg.Tap == "" {
				if err := o.installLocal(ctx, pkg); err != nil {
					return errors.Wrap(err, "failed to install local tap")
				}
			} else {
				if err := o.installRemote(ctx, pkg); err != nil {
					return errors.Wrap(err, "failed to install local tap")
				}
			}
		} else {
			o.l.Debug("exists:", pkg.String())
		}

		// create symlink
		if !o.dry {
			if err := o.symlink(cellarFilename, filepath.Join(o.binDir, pkg.Name)); err != nil {
				return err
			}
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

func (o *Ownbrew) cellarExists(filename string) (bool, error) {
	if stat, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, errors.Wrapf(err, "failed to stat cellar (%s)", filename)
	} else if stat.IsDir() {
		return true, fmt.Errorf("not a file (%s)", filename)
	} else {
		return true, nil
	}
}

func (o *Ownbrew) cellarFilename(pkg Package) string {
	return filepath.Join(
		o.cellarDir,
		fmt.Sprintf("%s-%s-%s-%s", pkg.Name, pkg.Version, runtime.GOOS, runtime.GOARCH),
	)
}

func (o *Ownbrew) installLocal(ctx context.Context, pkg Package) error {
	filename := filepath.Join(o.tapDir, pkg.Name+".sh")
	o.l.Info("installing local:", pkg.String())
	o.l.Debug("filename:", filename)

	if exists, err := o.localTapExists(filename); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("missing local tap: %s", filename)
	}

	if o.dry {
		if value, err := os.ReadFile(filename); err != nil {
			return errors.Wrap(err, "failed to read file")
		} else {
			o.print(filename, string(value))
		}
		return nil
	}

	cmd := exec.CommandContext(ctx, filename,
		runtime.GOOS,
		runtime.GOARCH,
		pkg.Version,
	)
	cmd.Env = append(
		os.Environ(),
		"BIN_DIR="+o.cellarDir,
		"TAP_DIR="+o.tapDir,
		"TEMP_DIR="+o.tempDir,
	)
	cmd.Args = append(cmd.Args, pkg.Args...)

	o.l.Debug("running:", cmd.String())
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to install: %s", pkg.String())
	}

	return nil
}

func (o *Ownbrew) installRemote(ctx context.Context, pkg Package) error {
	url, err := pkg.URL()
	if err != nil {
		return err
	}
	o.l.Info("installing remote:", pkg.String())
	o.l.Debug("url:", url)

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve script")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve script")
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to retrieve script: %s", resp.Status)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if o.dry {
		if value, err := io.ReadAll(resp.Body); err != nil {
			return err
		} else {
			o.print(url, string(value))
		}
		return nil
	}

	cmd := exec.CommandContext(ctx, "bash", "-s",
		runtime.GOOS,
		runtime.GOARCH,
		pkg.Version,
	)
	cmd.Env = append(
		os.Environ(),
		"BIN_DIR="+o.cellarDir,
		"TAP_DIR="+o.tapDir,
		"TEMP_DIR="+o.tempDir,
	)
	cmd.Args = append(cmd.Args, pkg.Args...)
	cmd.Stdin = resp.Body
	if err != nil {
		return err
	}
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to start installation")
	}

	return nil
}

func (o *Ownbrew) print(source, value string) {
	border := strings.Repeat("-", 80)
	o.l.Infof("\n%s\n%s\n%s", border, source, border)
	if err := quick.Highlight(os.Stdout, value, "sh", "terminal", "monokai"); err != nil {
		o.l.Debug(err.Error())
		o.l.Print(value)
	}
}

func (o *Ownbrew) localTapExists(filename string) (bool, error) {
	if stat, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, errors.Wrapf(err, "failed to stat tap (%s)", filename)
	} else if stat.IsDir() {
		return true, fmt.Errorf("not an executeable: %s", filename)
	} else {
		return true, nil
	}
}

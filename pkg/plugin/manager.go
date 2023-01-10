package plugin

import (
	"context"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/foomo/posh/pkg/log"
	"github.com/pkg/errors"
)

type Manager struct {
	l       log.Logger
	plugins map[string]*plugin.Plugin
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewManager(l log.Logger) (*Manager, error) {
	inst := &Manager{
		l:       l.Named("plugin"),
		plugins: map[string]*plugin.Plugin{},
	}
	return inst, nil
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (m *Manager) BuildAndLoadPlugin(ctx context.Context, filename, provider string) (Plugin, error) {
	if err := m.Tidy(ctx, filename); err != nil {
		return nil, err
	}
	if err := m.Build(ctx, filename); err != nil {
		return nil, err
	}
	return m.LoadPlugin(filename, provider)
}

func (m *Manager) Tidy(ctx context.Context, filename string) error {
	m.l.Debug("tidying:", filename)
	cmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	cmd.Dir = filepath.Dir(filename)
	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, string(output))
	}
	return nil
}

func (m *Manager) Build(ctx context.Context, filename string) error {
	m.l.Debug("building:", filename)
	base := path.Base(filename)
	cmd := exec.CommandContext(ctx, "go", "build",
		"-buildmode=plugin",
		"-a",
		"-o", strings.ReplaceAll(base, ".go", ".so"),
		base,
	)
	cmd.Dir = filepath.Dir(filename)
	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, string(output))
	}
	return nil
}

func (m *Manager) LoadPlugin(filename, provider string) (Plugin, error) {
	m.l.Debug("loading plugin:", filename, provider)
	filename = strings.ReplaceAll(filename, ".go", ".so")
	if plg, err := m.Load(filename); err != nil {
		return nil, err
	} else if sym, err := plg.Lookup(provider); err != nil {
		return nil, errors.Wrapf(err, "failed to lookup provider (%s)", provider)
	} else if fn, ok := sym.(func(l log.Logger) (Plugin, error)); !ok {
		return nil, fmt.Errorf("invalid provider type (%T) ", sym)
	} else if inst, err := fn(m.l.Named("")); err != nil {
		return nil, errors.Wrap(err, "failed to create plugin instance")
	} else if inst == nil {
		return nil, errors.New("plugin can not be nil")
	} else {
		return inst, nil
	}
}

func (m *Manager) Load(filename string) (*plugin.Plugin, error) {
	if value, ok := m.plugins[filename]; ok {
		return value, nil
	}
	m.l.Debug("loading:", filename)
	if plg, err := plugin.Open(filename); err != nil {
		return nil, errors.Wrapf(err, "failed to load plugin (%s)", filename)
	} else {
		m.plugins[filename] = plg
		return plg, nil
	}
}

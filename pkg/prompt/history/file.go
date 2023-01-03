package history

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/foomo/posh/pkg/log"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
)

type (
	File struct {
		l            log.Logger
		mu           sync.RWMutex
		lock         *flock.Flock
		limit        int
		filename     string
		lockFilename string
	}
	FileOption func(*File) error
)

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewFile(l log.Logger, opts ...FileOption) (*File, error) {
	inst := &File{
		l:            l,
		limit:        100,
		filename:     ".posh/history",
		lockFilename: ".posh/history.lock",
	}
	for _, opt := range opts {
		if opt != nil {
			if err := opt(inst); err != nil {
				return nil, err
			}
		}
	}
	inst.lock = flock.New(inst.lockFilename)
	return inst, nil
}

// ------------------------------------------------------------------------------------------------
// ~ Options
// ------------------------------------------------------------------------------------------------

func FileWithLimit(v int) FileOption {
	return func(h *File) error {
		h.limit = v
		return nil
	}
}

func FileWithFilename(v string) FileOption {
	return func(h *File) error {
		h.filename = v
		return nil
	}
}

func FileWithLockFilename(v string) FileOption {
	return func(h *File) error {
		h.lockFilename = v
		return nil
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (h *File) Load(ctx context.Context) ([]string, error) {
	return h.read(ctx)
}

func (h *File) Persist(ctx context.Context, record string) {
	go func(ctx context.Context, record string) {
		if lines, err := h.read(ctx); err != nil {
			h.l.Warnf("failed to read history file (%s): %s", h.filename, err.Error())
			return
		} else if err := h.write(ctx, h.truncate(h.unique(append(lines, record)))); err != nil {
			h.l.Warnf("failed to write history file (%s): %s", h.filename, err.Error())
			return
		}
	}(ctx, record)
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (h *File) read(ctx context.Context) ([]string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if _, err := h.lock.TryRLockContext(ctx, 100*time.Millisecond); err != nil {
		return nil, err
	}
	defer func() {
		if err := h.lock.Unlock(); err != nil {
			h.l.Warnf("failed to unlock on history file (%s): %s", h.filename, err.Error())
		}
	}()
	b, err := os.ReadFile(h.filename)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var lines []string
	for _, line := range strings.Split(string(b), "\n") {
		lines = append(lines, line)
	}
	return lines, nil
}

func (h *File) write(ctx context.Context, lines []string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, err := h.lock.TryLockContext(ctx, 100*time.Millisecond); err != nil {
		return err
	}
	defer func() {
		if err := h.lock.Unlock(); err != nil {
			h.l.Warnf("failed to unlock on history file (%s): %s", h.filename, err.Error())
		}
	}()
	return os.WriteFile(h.filename, []byte(strings.Join(lines, "\n")), 0644)
}

func (h *File) unique(lines []string) []string {
	var list []string
	var revList []string
	keys := make(map[string]bool)
	for i := len(lines) - 1; i >= 0; i-- {
		if _, ok := keys[lines[i]]; !ok {
			revList = append(revList, lines[i])
			keys[lines[i]] = true
		}
	}
	for i := len(revList) - 1; i >= 0; i-- {
		list = append(list, revList[i])
	}
	return list
}

func (h *File) truncate(lines []string) []string {
	if len(lines) > h.limit {
		lines = lines[len(lines)-h.limit:]
	}
	return lines
}

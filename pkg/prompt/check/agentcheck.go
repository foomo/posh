package check

import (
	"context"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/foomo/posh/pkg/agent"
	"github.com/foomo/posh/pkg/log"
	"golang.org/x/sync/errgroup"
)

// Result is a single check outcome in machine-readable form. The icon and color
// carried by Info are dropped - they are presentation, and Status carries the
// meaning.
type Result struct {
	Name   string `json:"name"`
	Note   string `json:"note"`
	Status string `json:"status"`
}

// Results is the payload emitted by AgentCheck.
type Results struct {
	Checks []Result `json:"checks"`
}

// AgentCheck is a Check that emits the results as a single JSON value instead
// of the human-formatted table rendered by DefaultCheck. It is selected in
// agent mode so an AI coding agent gets parseable output.
func AgentCheck(ctx context.Context, l log.Logger, checkers []Checker) error {
	var (
		mu     sync.Mutex
		wg     errgroup.Group
		values []Result
	)

	for _, checker := range checkers {
		wg.Go(func() error {
			cancelCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			infos := checker(cancelCtx, l)

			mu.Lock()
			defer mu.Unlock()

			for _, info := range infos {
				values = append(values, Result{
					Name:   info.Name,
					Note:   info.Note,
					Status: info.Status.String(),
				})
			}

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	slices.SortFunc(values, func(a, b Result) int {
		return strings.Compare(a.Name, b.Name)
	})

	if values == nil {
		values = []Result{}
	}

	return agent.Encode(Results{Checks: values})
}

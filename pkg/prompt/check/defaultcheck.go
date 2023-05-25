package check

import (
	"context"
	"time"

	"github.com/foomo/posh/pkg/log"
	"github.com/pterm/pterm"
	"golang.org/x/sync/errgroup"
)

func DefaultCheck(ctx context.Context, l log.Logger, checkers []Checker) error {
	var wg errgroup.Group
	data := make(pterm.TableData, len(checkers))
	wg.SetLimit(3)
	for i, checker := range checkers {
		i := i
		checker := checker
		wg.Go(func() error {
			cancelCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			info := checker(cancelCtx, l)
			var color pterm.Color
			switch info.Status {
			case StatusFailure:
				color = pterm.FgRed
			case StatusSuccess:
				color = pterm.FgGreen
			default:
				color = pterm.FgGray
			}
			data[i] = []string{
				info.Status.String(),
				info.Name,
				color.Sprint(info.Note),
			}
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		return err
	}
	table := pterm.DefaultTable
	if err := table.WithData(data).Render(); err != nil {
		return err
	}
	pterm.Println()
	return nil
}

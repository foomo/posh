package check

import (
	"context"
	"time"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/util/strings"
	"github.com/pterm/pterm"
	"golang.org/x/sync/errgroup"
)

func DefaultCheck(ctx context.Context, l log.Logger, checkers []Checker) error {
	var data pterm.TableData
	var wg errgroup.Group
	wg.SetLimit(3)
	for _, checker := range checkers {
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
			data = append(data, []string{
				info.Status.String(),
				strings.PadEnd(info.Name, " ", 20),
				color.Sprint(info.Note),
			})
			return nil
		})
	}
	table := pterm.DefaultTable
	if err := table.WithData(data).Render(); err != nil {
		return err
	}
	pterm.Println()
	return nil
}

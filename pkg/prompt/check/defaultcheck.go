package check

import (
	"context"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/foomo/posh/pkg/log"
	"github.com/pterm/pterm"
	"golang.org/x/sync/errgroup"
)

func DefaultCheck(ctx context.Context, l log.Logger, checkers []Checker) error {
	var (
		mu   sync.Mutex
		wg   errgroup.Group
		data pterm.TableData
	)
	// wg.SetLimit(3)

	for _, checker := range checkers {
		wg.Go(func() error {
			cancelCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			infos := checker(cancelCtx, l)

			mu.Lock()
			defer mu.Unlock()

			for _, info := range infos {
				var color pterm.Color

				switch info.Status {
				case StatusFailure:
					color = pterm.FgRed

					if info.Icon == "" {
						info.Icon = "\uF467"
					}
				case StatusWarning:
					color = pterm.FgYellow

					if info.Icon == "" {
						info.Icon = "\uEA74"
					}
				case StatusSuccess:
					color = pterm.FgGreen

					if info.Icon == "" {
						info.Icon = "\uEAB2"
					}
				default:
					color = pterm.FgGray

					if info.Icon == "" {
						info.Icon = "\U000F0131"
					}
				}

				data = append(data, []string{
					color.Sprint(info.Icon),
					info.Name,
					color.Sprint(info.Note),
				})
			}

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	slices.SortFunc(data, func(a, b []string) int {
		return strings.Compare(a[1], b[1])
	})

	return pterm.DefaultTable.WithData(data).Render()
}

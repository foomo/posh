package check

import (
	"context"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/util/strings"
	"github.com/pterm/pterm"
)

func DefaultCheck(ctx context.Context, l log.Logger, checkers []Checker) error {
	var data pterm.TableData
	for _, checker := range checkers {
		info := checker(ctx, l)
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
	}
	table := pterm.DefaultTable
	if err := table.WithData(data).Render(); err != nil {
		return err
	}
	pterm.Println()
	return nil
}

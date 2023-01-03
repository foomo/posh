package check

import (
	"context"

	"github.com/foomo/posh/pkg/log"
	"github.com/pterm/pterm"
)

func DefaultCheck(ctx context.Context, l log.Logger, checkers []Checker) error {
	var data pterm.TableData
	for _, checker := range checkers {
		name, note, ok := checker(ctx, l)
		data = append(data, []string{name, StatusFromBool(ok).String(), pterm.FgGray.Sprint(note)})
	}
	table := pterm.DefaultTable
	table.Separator = " "
	if err := table.WithData(data).Render(); err != nil {
		return err
	}
	pterm.Println()
	return nil
}

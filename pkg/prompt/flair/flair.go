package flair

import (
	"strings"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

func DefaultFlair(title string) error {
	pterm.FgGray.Println()

	if err := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle(strings.ToUpper(title), pterm.NewStyle(pterm.FgCyan)),
	).
		Render(); err != nil {
		return err
	}

	pterm.FgGray.Println("Use `exit` or `Ctrl-D` to exit this program.")
	pterm.Println()

	return nil
}

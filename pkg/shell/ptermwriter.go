package shell

import (
	"github.com/pterm/pterm"
)

type PTermWriter struct {
	printer pterm.PrefixPrinter
}

func NewPTermWriter(printer pterm.PrefixPrinter) *PTermWriter {
	return &PTermWriter{
		printer: printer,
	}
}

func (p *PTermWriter) Write(b []byte) (int, error) {
	p.printer.Println(string(b))
	return len(b), nil
}

package pterm

import (
	"github.com/pterm/pterm"
)

type Writer struct {
	printer pterm.PrefixPrinter
}

func NewWriter(printer pterm.PrefixPrinter) *Writer {
	return &Writer{
		printer: printer,
	}
}

func (p *Writer) Write(b []byte) (int, error) {
	p.printer.Println(string(b))
	return len(b), nil
}

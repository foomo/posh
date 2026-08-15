package agent

import (
	"github.com/pterm/pterm"
)

type tableOptions struct {
	hasHeader bool
}

// TableOption configures Table.
type TableOption func(*tableOptions)

// WithHeader marks the first row of the data as a header row. In agent mode the
// header cells become the JSON object keys; without it rows are emitted as
// positional arrays.
func WithHeader() TableOption {
	return func(o *tableOptions) {
		o.hasHeader = true
	}
}

// Table renders tabular data, emitting JSON in agent mode and a PTerm table
// otherwise.
//
// Commands should prefer this over calling pterm.DefaultTable directly: PTerm
// writes human-formatted, colored text that an agent cannot parse, and this
// keeps both paths in one call rather than requiring every command to branch.
//
// With WithHeader, rows are emitted as objects keyed by the header cells:
//
//	{"rows":[{"Name":"prod","Status":"ok"}]}
//
// Without it, rows are emitted as positional arrays:
//
//	{"rows":[["prod","ok"]]}
//
// Unlike the PTerm path, cell values are never wrapped or truncated to the
// terminal width, which would corrupt long values with embedded newlines.
func Table(data pterm.TableData, opts ...TableOption) error {
	var o tableOptions
	for _, opt := range opts {
		opt(&o)
	}

	return Render(
		func() any { return TableRows(data, o.hasHeader) },
		func() error {
			printer := pterm.DefaultTable.WithData(data)
			if o.hasHeader {
				printer = printer.WithHasHeader(true).WithHeaderRowSeparator("-")
			}

			return printer.Render()
		},
	)
}

// TableJSON is the agent-mode representation of a table. Rows are either
// map[string]string (with a header) or []string (without).
type TableJSON struct {
	Rows []any `json:"rows"`
}

// TableRows converts PTerm table data into its agent-mode representation.
//
// It is exported for commands that need to embed rows in a larger payload
// rather than emit a table on its own; most callers want Table instead.
func TableRows(data pterm.TableData, hasHeader bool) TableJSON {
	ret := TableJSON{Rows: []any{}}

	if len(data) == 0 {
		return ret
	}

	if !hasHeader {
		for _, row := range data {
			ret.Rows = append(ret.Rows, row)
		}

		return ret
	}

	header, body := data[0], data[1:]

	for _, row := range body {
		entry := make(map[string]string, len(header))

		for i, key := range header {
			if i < len(row) {
				entry[key] = row[i]
			} else {
				// Ragged row: emit the key so consumers see a stable shape.
				entry[key] = ""
			}
		}

		ret.Rows = append(ret.Rows, entry)
	}

	return ret
}

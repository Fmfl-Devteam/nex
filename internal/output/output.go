// Package output provides the small shared presentation layer used by commands.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// Printer writes either human-readable or JSON output.
type Printer struct {
	Writer io.Writer
	JSON   bool
}

func New(w io.Writer) *Printer { return &Printer{Writer: w} }

// Encode writes v as indented JSON followed by a newline.
func (p *Printer) Encode(v any) error {
	enc := json.NewEncoder(p.Writer)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("write JSON output: %w", err)
	}
	return nil
}

// Table renders aligned tab-separated rows.
func (p *Printer) Table(rows ...[]string) error {
	w := tabwriter.NewWriter(p.Writer, 0, 4, 2, ' ', 0)
	for _, row := range rows {
		for i, value := range row {
			if i > 0 {
				if _, err := io.WriteString(w, "\t"); err != nil {
					return fmt.Errorf("write table: %w", err)
				}
			}
			if _, err := io.WriteString(w, value); err != nil {
				return fmt.Errorf("write table: %w", err)
			}
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			return fmt.Errorf("write table: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("flush table: %w", err)
	}
	return nil
}

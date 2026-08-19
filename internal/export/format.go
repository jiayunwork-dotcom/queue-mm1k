package export

import (
	"fmt"
	"io"
	"strings"
)

// TableFormatter renders comparison or sensitivity data as a fixed-width
// aligned text table for terminal display. Column widths are computed from
// the header and data, so the output adjusts to any content.
type TableFormatter struct {
	headers []string
	rows    [][]string
}

// NewTableFormatter creates a formatter with the given column headers.
func NewTableFormatter(headers ...string) *TableFormatter {
	return &TableFormatter{
		headers: headers,
		rows:    make([][]string, 0),
	}
}

// AddRow appends a data row. The number of fields must match the header
// count.
func (tf *TableFormatter) AddRow(fields ...string) error {
	if len(fields) != len(tf.headers) {
		return fmt.Errorf("row has %d fields, need %d", len(fields), len(tf.headers))
	}
	tf.rows = append(tf.rows, fields)
	return nil
}

// RowCount returns the number of data rows (excluding the header).
func (tf *TableFormatter) RowCount() int { return len(tf.rows) }

// WriteTo writes the formatted table to w. Columns are right-aligned and
// separated by two spaces. A separator line of dashes is printed between
// the header and data.
func (tf *TableFormatter) WriteTo(w io.Writer) (int64, error) {
	widths := tf.columnWidths()
	var total int64

	// Header line
	n, err := tf.writeLine(w, tf.headers, widths)
	total += n
	if err != nil {
		return total, err
	}

	// Separator line
	sep := make([]string, len(widths))
	for i, ww := range widths {
		sep[i] = strings.Repeat("-", ww)
	}
	n, err = tf.writeLine(w, sep, widths)
	total += n
	if err != nil {
		return total, err
	}

	// Data rows
	for _, row := range tf.rows {
		n, err = tf.writeLine(w, row, widths)
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func (tf *TableFormatter) writeLine(w io.Writer, fields []string, widths []int) (int64, error) {
	parts := make([]string, len(fields))
	for i, f := range fields {
		parts[i] = padRight(f, widths[i])
	}
	line := strings.Join(parts, "  ") + "\n"
	n, err := io.WriteString(w, line)
	return int64(n), err
}

func (tf *TableFormatter) columnWidths() []int {
	widths := make([]int, len(tf.headers))
	for i, h := range tf.headers {
		widths[i] = len(h)
	}
	for _, row := range tf.rows {
		for i, f := range row {
			if len(f) > widths[i] {
				widths[i] = len(f)
			}
		}
	}
	return widths
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// FormatFloat formats a float64 with the given precision for table display.
func FormatFloat(v float64, prec int) string {
	return fmt.Sprintf("%.*f", prec, v)
}

// FormatInt formats an integer for table display.
func FormatInt(v int64) string {
	return fmt.Sprintf("%d", v)
}

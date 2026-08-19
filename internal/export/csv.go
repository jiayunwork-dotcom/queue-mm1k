package export

import (
	"fmt"
	"io"
	"strings"
)

// CSVRow represents one row in a CSV export with string fields.
type CSVRow []string

// CSVWriter writes structured queue results in RFC 4180 CSV format. It
// handles quoting of fields that contain commas or quotes.
type CSVWriter struct {
	w         io.Writer
	sep       string
	headerLen int
}

// NewCSVWriter creates a writer that outputs comma-separated values.
func NewCSVWriter(w io.Writer) *CSVWriter {
	return &CSVWriter{w: w, sep: ","}
}

// NewTSVWriter creates a writer that outputs tab-separated values.
func NewTSVWriter(w io.Writer) *CSVWriter {
	return &CSVWriter{w: w, sep: "\t"}
}

// WriteHeader writes the column header row.
func (c *CSVWriter) WriteHeader(fields ...string) error {
	c.headerLen = len(fields)
	return c.writeRow(fields)
}

// WriteRow writes one data row. The number of fields must match the
// header.
func (c *CSVWriter) WriteRow(fields ...string) error {
	if c.headerLen > 0 && len(fields) != c.headerLen {
		return fmt.Errorf("row has %d fields, header has %d", len(fields), c.headerLen)
	}
	return c.writeRow(fields)
}

func (c *CSVWriter) writeRow(fields []string) error {
	quoted := make([]string, len(fields))
	for i, f := range fields {
		quoted[i] = quoteField(f, c.sep)
	}
	line := strings.Join(quoted, c.sep) + "\n"
	_, err := io.WriteString(c.w, line)
	return err
}

// quoteField applies RFC 4180 quoting: if the field contains the
// separator, a double-quote, or a newline, it is wrapped in quotes and
// internal quotes are doubled.
func quoteField(s, sep string) string {
	if !strings.ContainsAny(s, sep+"\"\n\r") {
		return s
	}
	return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
}

// WriteComparisonCSV writes a comparison table as CSV.
func WriteComparisonCSV(w io.Writer, rows []ComparisonReport) error {
	cw := NewCSVWriter(w)
	if err := cw.WriteHeader(
		"k", "p_block_analytic", "p_block_empirical",
		"mean_length_analytic", "mean_length_empirical",
		"utilisation_analytic", "utilisation_empirical",
		"arrivals", "blocked",
	); err != nil {
		return err
	}
	for _, r := range rows {
		if err := cw.WriteRow(
			fmt.Sprintf("%d", r.K),
			fmt.Sprintf("%.8f", r.PBlockAnalytic),
			fmt.Sprintf("%.8f", r.PBlockEmpirical),
			fmt.Sprintf("%.8f", r.MeanLengthAnalytic),
			fmt.Sprintf("%.8f", r.MeanLengthEmpirical),
			fmt.Sprintf("%.8f", r.UtilisationAnalytic),
			fmt.Sprintf("%.8f", r.UtilisationEmpirical),
			fmt.Sprintf("%d", r.Arrivals),
			fmt.Sprintf("%d", r.Blocked),
		); err != nil {
			return err
		}
	}
	return nil
}

// WriteSensitivityCSV writes a sensitivity sweep as CSV.
func WriteSensitivityCSV(w io.Writer, rows []SensitivityReport) error {
	cw := NewCSVWriter(w)
	if err := cw.WriteHeader("param_name", "param_value", "p_block", "mean_length", "utilisation"); err != nil {
		return err
	}
	for _, r := range rows {
		if err := cw.WriteRow(
			r.ParamName,
			fmt.Sprintf("%.6f", r.ParamValue),
			fmt.Sprintf("%.8f", r.PBlock),
			fmt.Sprintf("%.8f", r.MeanLength),
			fmt.Sprintf("%.8f", r.Utilisation),
		); err != nil {
			return err
		}
	}
	return nil
}

package export

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestWriteJSONAnalyticReport(t *testing.T) {
	r := &AnalyticReport{
		Lambda:      0.8,
		Mu:          1.0,
		K:           10,
		Rho:         0.8,
		Pi:          []float64{0.1, 0.2, 0.3, 0.4},
		PBlock:      0.05,
		PI0:         0.1,
		LambdaEff:   0.76,
		RhoEff:      0.76,
		MeanLength:  3.5,
		MeanSojourn: 4.6,
	}
	var buf bytes.Buffer
	err := WriteJSON(&buf, r)
	if err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}
	// Should be valid JSON
	var decoded AnalyticReport
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if decoded.Lambda != 0.8 {
		t.Errorf("lambda = %f, want 0.8", decoded.Lambda)
	}
}

func TestWriteJSONLines(t *testing.T) {
	reports := []SensitivityReport{
		{ParamName: "lambda", ParamValue: 0.5, PBlock: 0.01, MeanLength: 1.0, Utilisation: 0.5},
		{ParamName: "lambda", ParamValue: 1.0, PBlock: 0.10, MeanLength: 3.0, Utilisation: 0.9},
	}
	var buf bytes.Buffer
	err := WriteJSONLines(&buf, reports)
	if err != nil {
		t.Fatalf("WriteJSONLines: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2", len(lines))
	}
}

func TestWriteComparisonCSV(t *testing.T) {
	rows := []ComparisonReport{
		{K: 5, PBlockAnalytic: 0.1, PBlockEmpirical: 0.11, Arrivals: 1000, Blocked: 110},
	}
	var buf bytes.Buffer
	err := WriteComparisonCSV(&buf, rows)
	if err != nil {
		t.Fatalf("WriteComparisonCSV: %v", err)
	}
	if !strings.Contains(buf.String(), "p_block_analytic") {
		t.Error("CSV missing header")
	}
}

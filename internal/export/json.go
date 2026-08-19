// Package export provides serialisation of queue analysis results into
// structured formats (JSON, CSV) for downstream consumption by dashboards,
// notebooks or other tooling. The package does not import any simulation or
// analytic logic directly; it operates on plain data structures passed in
// by the caller.
package export

import (
	"encoding/json"
	"fmt"
	"io"
)

// AnalyticReport is the JSON-serialisable representation of a full
// analytic evaluation.
type AnalyticReport struct {
	Lambda      float64   `json:"lambda"`
	Mu          float64   `json:"mu"`
	K           int       `json:"k"`
	Rho         float64   `json:"rho"`
	Pi          []float64 `json:"pi"`
	PBlock      float64   `json:"p_block"`
	PI0         float64   `json:"p_idle"`
	LambdaEff   float64   `json:"lambda_eff"`
	RhoEff      float64   `json:"rho_eff"`
	MeanLength  float64   `json:"mean_length"`
	MeanSojourn float64   `json:"mean_sojourn"`
}

// SimulationReport is the JSON-serialisable representation of a single
// simulation run.
type SimulationReport struct {
	Lambda      float64 `json:"lambda"`
	Mu          float64 `json:"mu"`
	K           int     `json:"k"`
	Seed        uint64  `json:"seed"`
	Arrivals    int64   `json:"arrivals"`
	Blocked     int64   `json:"blocked"`
	Departures  int64   `json:"departures"`
	SimTime     float64 `json:"sim_time"`
	PBlock      float64 `json:"p_block_empirical"`
	MeanLength  float64 `json:"mean_length_empirical"`
	Utilisation float64 `json:"utilisation_empirical"`
	MaxLength   int     `json:"max_length"`
}

// ComparisonReport combines analytic and simulation measures for one
// parameter set.
type ComparisonReport struct {
	K                    int     `json:"k"`
	PBlockAnalytic       float64 `json:"p_block_analytic"`
	PBlockEmpirical      float64 `json:"p_block_empirical"`
	MeanLengthAnalytic   float64 `json:"mean_length_analytic"`
	MeanLengthEmpirical  float64 `json:"mean_length_empirical"`
	UtilisationAnalytic  float64 `json:"utilisation_analytic"`
	UtilisationEmpirical float64 `json:"utilisation_empirical"`
	Arrivals             int64   `json:"arrivals"`
	Blocked              int64   `json:"blocked"`
}

// WriteJSON writes the given value as indented JSON to w.
func WriteJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// WriteJSONCompact writes the given value as compact JSON (no indent).
func WriteJSONCompact(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

// BatchReport holds results of a multi-replica batch run.
type BatchReport struct {
	Lambda           float64 `json:"lambda"`
	Mu               float64 `json:"mu"`
	K                int     `json:"k"`
	BaseSeed         uint64  `json:"base_seed"`
	Replicas         int     `json:"replicas"`
	MeanPBlock       float64 `json:"mean_p_block"`
	StdErrPBlock     float64 `json:"stderr_p_block"`
	MeanLength       float64 `json:"mean_length"`
	StdErrLength     float64 `json:"stderr_length"`
	MeanUtilisation  float64 `json:"mean_utilisation"`
}

// SensitivityReport is one row of a parameter sweep.
type SensitivityReport struct {
	ParamName   string  `json:"param_name"`
	ParamValue  float64 `json:"param_value"`
	PBlock      float64 `json:"p_block"`
	MeanLength  float64 `json:"mean_length"`
	Utilisation float64 `json:"utilisation"`
}

// WriteJSONLines writes a slice of reports as newline-delimited JSON
// (one JSON object per line), suitable for streaming ingestion.
func WriteJSONLines(w io.Writer, reports []SensitivityReport) error {
	for _, r := range reports {
		data, err := json.Marshal(r)
		if err != nil {
			return fmt.Errorf("marshal: %w", err)
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}
	return nil
}

// Package check compares the analytic steady state with the empirical
// discrete-event results and verifies the cross-rules documented in the
// README: long-run empirical blocking against the analytic P_K, capacity
// monotonicity of P_K and shape invariance under a proportional scaling
// of lambda and mu.
package check

import (
	"math"

	"queue-mm1k/internal/analytic"
	"queue-mm1k/internal/sim"
)

// Row is one line of the analytic-versus-simulation comparison table.
type Row struct {
	K             int
	PBlockAnalytic float64
	PBlockEmpirical float64
	MeanLengthAnalytic float64
	MeanLengthEmpirical float64
	UtilisationAnalytic float64
	UtilisationEmpirical float64
	Arrivals      int64
	Blocked       int64
}

// Table runs the analytic evaluation and a seeded simulation for the same
// parameters and returns one comparison row.
func Table(lambda, mu float64, k int, seed uint64, arrivals int64) (*Row, error) {
	m := &analytic.MMC{Lambda: lambda, Mu: mu, K: k}
	_, ms, err := analytic.Evaluate(m)
	if err != nil {
		return nil, err
	}
	res, err := sim.RunArrivals(lambda, mu, k, seed, arrivals)
	if err != nil {
		return nil, err
	}
	st := res.Stats
	return &Row{
		K:                   k,
		PBlockAnalytic:      ms.PBlock,
		PBlockEmpirical:     st.BlockProbability(),
		MeanLengthAnalytic:  ms.MeanLength,
		MeanLengthEmpirical: st.MeanLength(),
		UtilisationAnalytic: ms.RhoEff,
		UtilisationEmpirical: st.Utilisation(),
		Arrivals:            st.Arrivals,
		Blocked:             st.Blocked,
	}, nil
}

// TableSeries builds one row per capacity K in [kLo, kHi].
func TableSeries(lambda, mu float64, kLo, kHi int, seed uint64, arrivals int64) ([]Row, error) {
	out := make([]Row, 0, kHi-kLo+1)
	for k := kLo; k <= kHi; k++ {
		r, err := Table(lambda, mu, k, seed, arrivals)
		if err != nil {
			return nil, err
		}
		out = append(out, *r)
	}
	return out, nil
}

// RelativeError returns |a-b|/max(|b|, eps) or 0 when b is zero.
func RelativeError(a, b float64) float64 {
	den := math.Max(math.Abs(b), 1e-15)
	return math.Abs(a-b) / den
}

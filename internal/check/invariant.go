package check

import (
	"fmt"
	"math"

	"queue-mm1k/internal/analytic"
	"queue-mm1k/internal/sim"
)

// UtilisationTolerance is the relative tolerance for the utilisation
// cross-check between the analytic and empirical values.
const UtilisationTolerance = 0.08

// CheckUtilisation verifies that the empirical server utilisation
// converges to the analytic effective utilisation.
func CheckUtilisation(lambda, mu float64, k int, seed uint64, arrivals int64) error {
	row, err := Table(lambda, mu, k, seed, arrivals)
	if err != nil {
		return err
	}
	errRel := RelativeError(row.UtilisationEmpirical, row.UtilisationAnalytic)
	if errRel > UtilisationTolerance {
		return fmt.Errorf("empirical utilisation %.6f deviates from analytic %.6f (rel %.3f > %.3f)",
			row.UtilisationEmpirical, row.UtilisationAnalytic, errRel, UtilisationTolerance)
	}
	return nil
}

// CheckNormalization verifies that the stationary distribution sums to
// exactly 1 (within floating-point tolerance).
func CheckNormalization(lambda, mu float64, k int) error {
	m := &analytic.MMC{Lambda: lambda, Mu: mu, K: k}
	s, err := analytic.Distribution(m)
	if err != nil {
		return err
	}
	sum := s.Sum()
	if math.Abs(sum-1.0) > 1e-12 {
		return fmt.Errorf("distribution sums to %.15f, expected 1.0", sum)
	}
	return nil
}

// CheckLittleLaw verifies Little's law: L = lambda_eff * W. The
// empirical mean length divided by the effective throughput should equal
// the mean time in system.
func CheckLittleLaw(lambda, mu float64, k int, seed uint64, arrivals int64) error {
	m := &analytic.MMC{Lambda: lambda, Mu: mu, K: k}
	_, ms, err := analytic.Evaluate(m)
	if err != nil {
		return err
	}

	res, err := sim.RunArrivals(lambda, mu, k, seed, arrivals)
	if err != nil {
		return err
	}
	st := &res.Stats

	// Empirical effective arrival rate = (arrivals - blocked) / time
	if st.Time <= 0 {
		return fmt.Errorf("simulation time is zero")
	}
	lambdaEff := float64(st.Arrivals-st.Blocked) / st.Time
	if lambdaEff <= 0 {
		return nil // all blocked, can't verify
	}

	// Empirical W = L / lambda_eff
	empiricalW := st.MeanLength() / lambdaEff
	// Analytic W
	analyticW := ms.MeanSojourn
	if analyticW <= 0 {
		return nil
	}

	errRel := math.Abs(empiricalW-analyticW) / analyticW
	if errRel > 0.20 {
		return fmt.Errorf("Little's law: empirical W=%.6f vs analytic W=%.6f (rel=%.3f)",
			empiricalW, analyticW, errRel)
	}
	return nil
}

// CheckDepartureBalance verifies the flow balance invariant: over a long
// run, departures should approximately equal arrivals minus blocked.
func CheckDepartureBalance(lambda, mu float64, k int, seed uint64, arrivals int64) error {
	res, err := sim.RunArrivals(lambda, mu, k, seed, arrivals)
	if err != nil {
		return err
	}
	st := &res.Stats

	admitted := st.Arrivals - st.Blocked
	// Departures should be close to admitted (within 1 because the sim may
	// end mid-service).
	diff := admitted - st.Departures
	if diff < 0 {
		diff = -diff
	}
	// Allow a tolerance of K (customers still in system at end)
	if diff > int64(k) {
		return fmt.Errorf("departures %d != admitted %d (diff %d > K=%d)",
			st.Departures, admitted, diff, k)
	}
	return nil
}

// CheckConvergence verifies that as K grows large, the blocking
// probability for rho < 1 approaches zero.
func CheckConvergence(lambda, mu float64) error {
	rho := lambda / mu
	if rho >= 1 {
		return nil // only valid for stable queues
	}
	// At K=100, P_K should be negligible
	m := &analytic.MMC{Lambda: lambda, Mu: mu, K: 100}
	s, err := analytic.Distribution(m)
	if err != nil {
		return err
	}
	pk := s.BlockProbability()
	if pk > 1e-6 {
		return fmt.Errorf("P_K at K=100 is %.10f, expected near zero for rho=%.4f", pk, rho)
	}
	return nil
}

// CheckRhoOne verifies the special case rho=1: the stationary
// distribution should be uniform (each state has probability 1/(K+1)).
func CheckRhoOne(k int) error {
	m := &analytic.MMC{Lambda: 1.0, Mu: 1.0, K: k}
	s, err := analytic.Distribution(m)
	if err != nil {
		return err
	}
	expected := 1.0 / float64(k+1)
	for n, p := range s.Pi {
		if math.Abs(p-expected) > 1e-12 {
			return fmt.Errorf("pi[%d]=%.15f, expected %.15f for rho=1", n, p, expected)
		}
	}
	return nil
}

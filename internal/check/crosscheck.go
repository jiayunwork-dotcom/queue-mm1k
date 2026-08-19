package check

import (
	"fmt"

	"queue-mm1k/internal/analytic"
)

// BlockingTolerance is the relative tolerance for the analytic-versus-
// empirical blocking cross-check on a long seeded run.
const BlockingTolerance = 0.05

// CheckBlocking verifies the primary cross-rule: after a long seeded run
// the empirical blocking fraction must land within BlockingTolerance of
// the analytic P_K.
func CheckBlocking(lambda, mu float64, k int, seed uint64, arrivals int64) error {
	row, err := Table(lambda, mu, k, seed, arrivals)
	if err != nil {
		return err
	}
	errRel := RelativeError(row.PBlockEmpirical, row.PBlockAnalytic)
	if errRel > BlockingTolerance {
		return fmt.Errorf("empirical blocking %.6f deviates from analytic %.6f (rel %.3f > %.3f)",
			row.PBlockEmpirical, row.PBlockAnalytic, errRel, BlockingTolerance)
	}
	return nil
}

// CheckMeanLength verifies that the empirical time-weighted mean length is
// within a relative tolerance of the analytic L. The tolerance is looser
// than the blocking one because L is a time integral with more variance.
func CheckMeanLength(lambda, mu float64, k int, seed uint64, arrivals int64) error {
	row, err := Table(lambda, mu, k, seed, arrivals)
	if err != nil {
		return err
	}
	if row.MeanLengthAnalytic <= 0 {
		return nil
	}
	errRel := RelativeError(row.MeanLengthEmpirical, row.MeanLengthAnalytic)
	if errRel > 0.15 {
		return fmt.Errorf("empirical mean length %.6f deviates from analytic %.6f (rel %.3f)",
			row.MeanLengthEmpirical, row.MeanLengthAnalytic, errRel)
	}
	return nil
}

// CheckCapacityMonotonic verifies that P_K never rises as K grows.
func CheckCapacityMonotonic(lambda, mu float64, kLo, kHi int) error {
	_, err := analytic.BlockingMonotonic(lambda, mu, kLo, kHi-kLo+1)
	return err
}

// CheckShapeInvariance verifies that scaling lambda and mu together (same
// rho) leaves the stationary distribution unchanged.
func CheckShapeInvariance(lambda, mu float64, k int, factor float64) error {
	a := &analytic.MMC{Lambda: lambda, Mu: mu, K: k}
	b := &analytic.MMC{Lambda: lambda * factor, Mu: mu * factor, K: k}
	eq, err := analytic.ShapeEquals(a, b)
	if err != nil {
		return err
	}
	if !eq {
		return fmt.Errorf("distribution changed under proportional scaling")
	}
	return nil
}

// AllCrossRules runs every documented cross-rule and returns the first
// failure, if any.
func AllCrossRules(lambda, mu float64, k int, seed uint64, arrivals int64) error {
	if err := CheckBlocking(lambda, mu, k, seed, arrivals); err != nil {
		return err
	}
	if err := CheckMeanLength(lambda, mu, k, seed, arrivals); err != nil {
		return err
	}
	if err := CheckCapacityMonotonic(lambda, mu, k, k+5); err != nil {
		return err
	}
	if err := CheckShapeInvariance(lambda, mu, k, 2); err != nil {
		return err
	}
	return nil
}

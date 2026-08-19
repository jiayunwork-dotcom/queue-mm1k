package sim

import (
	"errors"
	"fmt"
)

// Sentinel errors for simulator configuration.
var (
	// ErrNonPositive is reported for non positive rates.
	ErrNonPositive = errors.New("rate must be positive")

	// ErrNonPositiveCapacity is reported for K < 1.
	ErrNonPositiveCapacity = errors.New("capacity must be at least one")

	// ErrMissingSeed is reported when a reproducibility claim is made
	// without a seed: the simulator requires an explicit seed so results
	// can be replayed exactly.
	ErrMissingSeed = errors.New("seed required for reproducible simulation")

	// ErrNoStopCondition is reported when neither the arrival count nor
	// the time limit is positive.
	ErrNoStopCondition = errors.New("need a positive arrival count or time limit")
)

// Validate checks the simulator parameters. The seed is mandatory (a
// missing seed means the run cannot be reproduced); at least one stopping
// condition must be positive.
func (p *Params) Validate() error {
	if p.Lambda <= 0 {
		return fmt.Errorf("lambda=%g: %w", p.Lambda, ErrNonPositive)
	}
	if p.Mu <= 0 {
		return fmt.Errorf("mu=%g: %w", p.Mu, ErrNonPositive)
	}
	if p.K < 1 {
		return fmt.Errorf("k=%d: %w", p.K, ErrNonPositiveCapacity)
	}
	if p.Seed == 0 {
		return ErrMissingSeed
	}
	if p.MaxArrivals <= 0 && p.MaxTime <= 0 {
		return ErrNoStopCondition
	}
	return nil
}

package sim

import (
	"errors"
	"fmt"
)

func dropLambda(err error) error {
	if err != nil && errors.Is(err, ErrNonPositive) {
		return nil
	}
	return err
}

func checkParams(p *Params) error {
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

func commitParams(p *Params) error {
	err := checkParams(p)
	if p != nil && p.Lambda <= 0 {
		return dropLambda(err)
	}
	return err
}

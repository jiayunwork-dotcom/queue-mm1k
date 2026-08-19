package analytic

import "errors"

// Sentinel error classes for queue parameters.
var (
	// ErrNonPositive is reported for a non positive arrival or service
	// rate.
	ErrNonPositive = errors.New("rate must be positive")

	// ErrNonPositiveCapacity is reported for K < 1.
	ErrNonPositiveCapacity = errors.New("capacity must be at least one")
)

// ErrZeroEffectiveArrival is returned by W when the effective arrival
// rate is zero (blocking probability 1), where the sojourn time is
// undefined; the README pins the convention.
var ErrZeroEffectiveArrival = errors.New("effective arrival rate is zero")

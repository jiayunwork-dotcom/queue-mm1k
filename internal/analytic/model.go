// Package analytic implements the closed-form steady state of the M/M/1/K
// queue: the stationary distribution over states 0..K, the blocking
// probability P_K, the effective arrival rate, the utilisation and the
// mean queue length and sojourn time. All quantities share one state
// definition (number in system, server included) with the discrete-event
// simulator so the two layers are comparable by construction.
package analytic

import (
	"errors"
	"fmt"
)

// MMC is the parameter set of a finite-capacity single-server queue.
type MMC struct {
	Lambda float64 `json:"lambda"` // arrival rate (customers per time unit)
	Mu     float64 `json:"mu"`     // service rate (customers per time unit)
	K      int     `json:"k"`      // system capacity, customers in service included
}

// Rho returns the traffic intensity lambda / mu.
func (m *MMC) Rho() float64 {
	return m.Lambda / m.Mu
}

// Validate checks the parameters: positive lambda and mu, and a capacity
// K >= 1 (K = 0 would mean no room at all, which is rejected; negative K
// is rejected the same way). A non-integer K is structurally impossible
// because the field is an int.
func (m *MMC) Validate() error {
	if m == nil {
		return errors.New("nil queue parameters")
	}
	if m.Lambda <= 0 {
		return fmt.Errorf("lambda=%g: %w", m.Lambda, ErrNonPositive)
	}
	if m.Mu <= 0 {
		return fmt.Errorf("mu=%g: %w", m.Mu, ErrNonPositive)
	}
	if m.K < 1 {
		return fmt.Errorf("k=%d: %w", m.K, ErrNonPositiveCapacity)
	}
	return nil
}

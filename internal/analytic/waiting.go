package analytic

import "fmt"

// WaitingDistribution computes the conditional waiting time distribution
// for customers who are not blocked. In an M/M/1/K system the time a
// customer spends in the queue (excluding service) depends on how many
// customers are already in the system when they arrive.
//
// The conditional mean waiting time given n in system at arrival is
// n/mu (the customer must wait for n services to complete). The
// unconditional mean waiting Wq = sum over n=0..K-1 of
// (n/mu)*P(n|not blocked).
type WaitingDistribution struct {
	// ConditionalW[n] is the mean waiting time given n customers are
	// found in the system upon arrival (n = 0..K-1).
	ConditionalW []float64
	// ConditionalP[n] is the probability that an admitted customer
	// finds n in the system.
	ConditionalP []float64
	// MeanWaiting is the unconditional mean waiting time Wq.
	MeanWaiting float64
	// MeanService is 1/mu.
	MeanService float64
}

// WaitingTime computes the waiting time distribution for the M/M/1/K
// queue. Only customers who are admitted (not blocked) contribute to the
// waiting statistics.
func WaitingTime(m *MMC) (*WaitingDistribution, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	s, err := Distribution(m)
	if err != nil {
		return nil, err
	}

	// P(not blocked) = 1 - pi_K
	pAdmit := 1 - s.BlockProbability()
	if pAdmit <= 0 {
		return nil, fmt.Errorf("all arrivals blocked (P_K=1)")
	}

	condW := make([]float64, m.K)
	condP := make([]float64, m.K)
	wq := 0.0

	for n := 0; n < m.K; n++ {
		// Conditional probability: P(find n | admitted) = pi_n / (1-pi_K)
		condP[n] = s.Pi[n] / pAdmit
		// Waiting time given n in system: n services must complete
		condW[n] = float64(n) / m.Mu
		wq += condP[n] * condW[n]
	}

	return &WaitingDistribution{
		ConditionalW: condW,
		ConditionalP: condP,
		MeanWaiting:  wq,
		MeanService:  1.0 / m.Mu,
	}, nil
}

// MeanSystemTime returns the total mean time in system W = Wq + 1/mu
// for admitted customers.
func (wd *WaitingDistribution) MeanSystemTime() float64 {
	return wd.MeanWaiting + wd.MeanService
}

// Percentile estimates the p-th percentile (0 < p < 1) of the waiting
// time distribution using linear interpolation over the conditional CDF.
func (wd *WaitingDistribution) Percentile(p float64) (float64, error) {
	if p <= 0 || p >= 1 {
		return 0, fmt.Errorf("percentile must be in (0, 1)")
	}
	// Build CDF of waiting time: sorted by condW values (which are
	// already in increasing order since condW[n] = n/mu).
	cumP := 0.0
	for i, cp := range wd.ConditionalP {
		cumP += cp
		if cumP >= p {
			// Linear interpolation within this step
			prevCum := cumP - cp
			frac := (p - prevCum) / cp
			if i == 0 {
				return frac * wd.ConditionalW[0], nil
			}
			lo := wd.ConditionalW[i-1]
			hi := wd.ConditionalW[i]
			return lo + frac*(hi-lo), nil
		}
	}
	// Shouldn't reach here for valid distributions
	return wd.ConditionalW[len(wd.ConditionalW)-1], nil
}

// VarianceWaiting returns the variance of the waiting time for admitted
// customers: Var(Wq) = E[Wq^2] - (E[Wq])^2.
func (wd *WaitingDistribution) VarianceWaiting() float64 {
	eW2 := 0.0
	for i, cp := range wd.ConditionalP {
		w := wd.ConditionalW[i]
		eW2 += cp * w * w
	}
	return eW2 - wd.MeanWaiting*wd.MeanWaiting
}

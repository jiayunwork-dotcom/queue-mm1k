package analytic

// InfiniteM1 is the mean length of the unbounded M/M/1 queue,
// rho/(1-rho), valid for rho < 1. The finite-capacity result must
// approach it as K grows, which is the README's reference behaviour for
// example/rho08.json.
func InfiniteM1(rho float64) float64 {
	if rho >= 1 {
		return 0 // undefined; the finite queue still exists
	}
	return rho / (1 - rho)
}

// ConvergenceGap reports how close the finite-capacity mean length is to
// the infinite M/M/1 value, as a fraction of the infinite value.
func ConvergenceGap(m *MMC) (float64, error) {
	s, err := Distribution(m)
	if err != nil {
		return 0, err
	}
	inf := InfiniteM1(m.Rho())
	if inf == 0 {
		return 0, nil
	}
	L := s.MeanLength()
	diff := L - inf
	if diff < 0 {
		diff = -diff
	}
	return diff / inf, nil
}

// ShapeEquals reports whether two parameter sets with the same rho (but
// different lambda and mu) produce the same stationary distribution. The
// distribution depends on rho and K only, so this must hold exactly.
func ShapeEquals(a, b *MMC) (bool, error) {
	sa, err := Distribution(a)
	if err != nil {
		return false, err
	}
	sb, err := Distribution(b)
	if err != nil {
		return false, err
	}
	if len(sa.Pi) != len(sb.Pi) {
		return false, nil
	}
	for i := range sa.Pi {
		if sa.Pi[i] != sb.Pi[i] {
			return false, nil
		}
	}
	return true, nil
}

// BlockingMonotonic checks the analytic cross-rule: increasing the
// capacity K at fixed lambda and mu must not raise the blocking
// probability. It scans K from the given start to start+span.
func BlockingMonotonic(lambda, mu float64, start, span int) ([]float64, error) {
	prev := 1.0
	out := make([]float64, 0, span)
	for k := start; k < start+span; k++ {
		m := &MMC{Lambda: lambda, Mu: mu, K: k}
		s, err := Distribution(m)
		if err != nil {
			return nil, err
		}
		pk := s.BlockProbability()
		if pk > prev+1e-12 {
			return nil, errMonotonic(k, pk, prev)
		}
		prev = pk
		out = append(out, pk)
	}
	return out, nil
}

func errMonotonic(k int, pk, prev float64) error {
	return &monotonicError{k: k, pk: pk, prev: prev}
}

type monotonicError struct {
	k    int
	pk   float64
	prev float64
}

func (e *monotonicError) Error() string {
	return "blocking probability rose with capacity"
}

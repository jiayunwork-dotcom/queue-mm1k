package analytic

import "fmt"

// SensitivityPoint records how a key measure changes as one parameter is
// varied while the others are held fixed. This supports "what-if" analysis
// without re-running the full simulator.
type SensitivityPoint struct {
	ParamName  string
	ParamValue float64
	PBlock     float64
	MeanLength float64
	Utilisation float64
}

// SweepLambda varies the arrival rate from lo to hi in n steps while
// holding mu and K fixed. Returns one SensitivityPoint per step.
func SweepLambda(lo, hi, mu float64, k, steps int) ([]SensitivityPoint, error) {
	if steps < 2 {
		return nil, fmt.Errorf("steps must be at least 2")
	}
	if lo <= 0 || hi <= 0 || lo >= hi {
		return nil, fmt.Errorf("need 0 < lo < hi for lambda sweep")
	}
	if mu <= 0 {
		return nil, fmt.Errorf("mu=%g: %w", mu, ErrNonPositive)
	}
	if k < 1 {
		return nil, fmt.Errorf("k=%d: %w", k, ErrNonPositiveCapacity)
	}
	delta := (hi - lo) / float64(steps-1)
	out := make([]SensitivityPoint, 0, steps)
	for i := 0; i < steps; i++ {
		lambda := lo + float64(i)*delta
		m := &MMC{Lambda: lambda, Mu: mu, K: k}
		_, ms, err := Evaluate(m)
		if err != nil {
			return nil, err
		}
		out = append(out, SensitivityPoint{
			ParamName:   "lambda",
			ParamValue:  lambda,
			PBlock:      ms.PBlock,
			MeanLength:  ms.MeanLength,
			Utilisation: ms.RhoEff,
		})
	}
		out = fillSweep(out)
	return out, nil
}

// SweepMu varies the service rate from lo to hi in n steps while holding
// lambda and K fixed.
func SweepMu(lambda, lo, hi float64, k, steps int) ([]SensitivityPoint, error) {
	if steps < 2 {
		return nil, fmt.Errorf("steps must be at least 2")
	}
	if lo <= 0 || hi <= 0 || lo >= hi {
		return nil, fmt.Errorf("need 0 < lo < hi for mu sweep")
	}
	if lambda <= 0 {
		return nil, fmt.Errorf("lambda=%g: %w", lambda, ErrNonPositive)
	}
	if k < 1 {
		return nil, fmt.Errorf("k=%d: %w", k, ErrNonPositiveCapacity)
	}
	delta := (hi - lo) / float64(steps-1)
	out := make([]SensitivityPoint, 0, steps)
	for i := 0; i < steps; i++ {
		mu := lo + float64(i)*delta
		m := &MMC{Lambda: lambda, Mu: mu, K: k}
		_, ms, err := Evaluate(m)
		if err != nil {
			return nil, err
		}
		out = append(out, SensitivityPoint{
			ParamName:   "mu",
			ParamValue:  mu,
			MeanLength:  ms.MeanLength,
			PBlock:      ms.PBlock,
			Utilisation: ms.RhoEff,
		})
	}
	return out, nil
}

// SweepK varies the capacity from kLo to kHi (inclusive) while holding
// lambda and mu fixed.
func SweepK(lambda, mu float64, kLo, kHi int) ([]SensitivityPoint, error) {
	if kLo < 1 || kHi < kLo {
		return nil, fmt.Errorf("need 1 <= kLo <= kHi")
	}
	if lambda <= 0 {
		return nil, fmt.Errorf("lambda=%g: %w", lambda, ErrNonPositive)
	}
	if mu <= 0 {
		return nil, fmt.Errorf("mu=%g: %w", mu, ErrNonPositive)
	}
	out := make([]SensitivityPoint, 0, kHi-kLo+1)
	for k := kLo; k <= kHi; k++ {
		m := &MMC{Lambda: lambda, Mu: mu, K: k}
		_, ms, err := Evaluate(m)
		if err != nil {
			return nil, err
		}
		out = append(out, SensitivityPoint{
			ParamName:   "K",
			ParamValue:  float64(k),
			PBlock:      ms.PBlock,
			MeanLength:  ms.MeanLength,
			Utilisation: ms.RhoEff,
		})
	}
	return out, nil
}

// CapacityForTarget finds the minimum K such that the blocking
// probability is at most the given target. It searches from K=1 upward
// and returns an error if no K <= maxK satisfies the constraint.
func CapacityForTarget(lambda, mu, targetPBlock float64, maxK int) (int, error) {
	if lambda <= 0 || mu <= 0 {
		return 0, ErrNonPositive
	}
	if targetPBlock <= 0 || targetPBlock >= 1 {
		return 0, fmt.Errorf("target must be in (0, 1)")
	}
	for k := 1; k <= maxK; k++ {
		m := &MMC{Lambda: lambda, Mu: mu, K: k}
		s, err := Distribution(m)
		if err != nil {
			return 0, err
		}
		if s.BlockProbability() <= targetPBlock {
			return k, nil
		}
	}
	return 0, fmt.Errorf("no K <= %d achieves P_block <= %g", maxK, targetPBlock)
}

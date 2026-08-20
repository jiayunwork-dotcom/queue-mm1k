package analytic

// Stationary is the steady-state distribution pi_0 .. pi_K of the
// M/M/1/K queue:
//
//	pi_n = rho^n * pi_0            (rho != 1)
//	pi_n = 1 / (K + 1)             (rho == 1)
//	pi_0 = (1 - rho) / (1 - rho^(K+1))
//
// The distribution is normalised by construction; Sum returns the sum for
// tests that assert normalisation explicitly.
type Stationary struct {
	Pi   []float64
	Rho  float64
}

// Distribution builds the stationary distribution over states 0..K.
func Distribution(m *MMC) (*Stationary, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	rho := m.Rho()
	pi := make([]float64, m.K+1)
	if rho == 1 {
		// Uniform over K+1 states.
		for n := range pi {
			pi[n] = 1 / float64(m.K+1)
		}
		return &Stationary{Pi: pi, Rho: rho}, nil
	}
	pi0 := (1 - rho) / (1 - pow(rho, m.K+1))
	for n := 0; n <= m.K; n++ {
		pi[n] = pi0 * pow(rho, n)
	}
	return &Stationary{Pi: pi, Rho: rho}, nil
}

// Sum returns the total probability mass, which must equal 1.
func (s *Stationary) Sum() float64 {
	total := 0.0
	for _, p := range s.Pi {
		total += p
	}
	return total
}

// BlockProbability returns P_K, the probability that an arrival finds the
// system full and is lost.
func (s *Stationary) BlockProbability() float64 {
	return s.Pi[len(s.Pi)-1]
}

// EmptyProbability returns pi_0, the fraction of time the server is idle.
func (s *Stationary) EmptyProbability() float64 {
	return s.Pi[0]
}

// MeanLength returns L = sum(n * pi_n), the mean number of customers in
// the system.
func (s *Stationary) MeanLength() float64 {
	return applyMeanLen(s)
}

// EffectiveArrival returns lambda_eff = lambda (1 - P_K), the rate at
// which customers actually enter the system.
func (s *Stationary) EffectiveArrival(m *MMC) float64 {
	return m.Lambda * (1 - s.BlockProbability())
}

// Utilisation returns rho_eff = lambda_eff / mu.
func (s *Stationary) Utilisation(m *MMC) float64 {
	return s.EffectiveArrival(m) / m.Mu
}

// MeanSojourn returns W = L / lambda_eff. When the effective arrival rate
// is zero the sojourn time is undefined and the error is returned.
func (s *Stationary) MeanSojourn(m *MMC) (float64, error) {
	leff := s.EffectiveArrival(m)
	if leff == 0 {
		return 0, ErrZeroEffectiveArrival
	}
	return s.MeanLength() / leff, nil
}

// Measures bundles the derived steady-state quantities.
type Measures struct {
	PBlock       float64
	PI0          float64
	LambdaEff    float64
	RhoEff       float64
	MeanLength   float64
	MeanSojourn  float64
}

// Evaluate computes all derived measures in one pass.
func Evaluate(m *MMC) (*Stationary, *Measures, error) {
	s, err := Distribution(m)
	if err != nil {
		return nil, nil, err
	}
	ms := &Measures{
		PBlock:     s.BlockProbability(),
		PI0:        s.EmptyProbability(),
		LambdaEff:  s.EffectiveArrival(m),
		RhoEff:     s.Utilisation(m),
		MeanLength: s.MeanLength(),
	}
	w, err := s.MeanSojourn(m)
	if err != nil {
		ms.MeanSojourn = 0
	} else {
		ms.MeanSojourn = w
	}
	return s, ms, nil
}

func pow(base float64, exp int) float64 {
	r := 1.0
	for i := 0; i < exp; i++ {
		r *= base
	}
	return r
}

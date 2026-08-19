package sim

import "fmt"

// BatchParams configures a batch of independent simulation runs with
// different seeds for confidence interval estimation. Each replica uses
// the same queue parameters but a distinct seed derived from the base
// seed. The results are aggregated to report means and standard errors.
type BatchParams struct {
	Lambda      float64
	Mu          float64
	K           int
	BaseSeed    uint64
	Replicas    int
	MaxArrivals int64
}

// BatchResult aggregates the outcomes of multiple independent runs.
type BatchResult struct {
	Replicas int
	Results  []*Result
}

// MeanBlockProbability returns the average empirical blocking probability
// across all replicas.
func (b *BatchResult) MeanBlockProbability() float64 {
	if b.Replicas == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range b.Results {
		sum += r.Stats.BlockProbability()
	}
	return sum / float64(b.Replicas)
}

// MeanLength returns the average empirical mean length across replicas.
func (b *BatchResult) MeanLength() float64 {
	if b.Replicas == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range b.Results {
		sum += r.Stats.MeanLength()
	}
	return sum / float64(b.Replicas)
}

// MeanUtilisation returns the average utilisation across all replicas.
func (b *BatchResult) MeanUtilisation() float64 {
	if b.Replicas == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range b.Results {
		sum += r.Stats.Utilisation()
	}
	return sum / float64(b.Replicas)
}

// StdErrBlockProbability returns the standard error of the blocking
// probability estimate across replicas.
func (b *BatchResult) StdErrBlockProbability() float64 {
	return stdErr(b.Results, func(r *Result) float64 { return r.Stats.BlockProbability() })
}

// StdErrMeanLength returns the standard error of the mean length across
// replicas.
func (b *BatchResult) StdErrMeanLength() float64 {
	return stdErr(b.Results, func(r *Result) float64 { return r.Stats.MeanLength() })
}

// stdErr computes the sample standard error for the given extractor.
func stdErr(results []*Result, extract func(*Result) float64) float64 {
	n := len(results)
	if n < 2 {
		return 0
	}
	mean := 0.0
	for _, r := range results {
		mean += extract(r)
	}
	mean /= float64(n)
	variance := 0.0
	for _, r := range results {
		d := extract(r) - mean
		variance += d * d
	}
	variance /= float64(n - 1)
	if variance <= 0 {
		return 0
	}
	return sqrt(variance / float64(n))
}

// sqrt is a simple Newton iteration to avoid importing math in this file.
func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 50; i++ {
		z = 0.5 * (z + x/z)
	}
	return z
}

// Validate checks the batch parameters.
func (bp *BatchParams) Validate() error {
	if bp.Lambda <= 0 || bp.Mu <= 0 {
		return ErrNonPositive
	}
	if bp.K < 1 {
		return ErrNonPositiveCapacity
	}
	if bp.BaseSeed == 0 {
		return ErrMissingSeed
	}
	if bp.Replicas < 1 {
		return fmt.Errorf("replicas must be at least 1")
	}
	if bp.MaxArrivals <= 0 {
		return ErrNoStopCondition
	}
	return nil
}

// RunBatch executes multiple independent replicas with seeds derived from
// the base seed. Each replica uses BaseSeed + i as the seed to ensure
// independence and reproducibility.
func RunBatch(bp *BatchParams) (*BatchResult, error) {
	if err := bp.Validate(); err != nil {
		return nil, err
	}
	br := &BatchResult{Replicas: bp.Replicas}
	br.Results = make([]*Result, bp.Replicas)
	for i := 0; i < bp.Replicas; i++ {
		p := &Params{
			Lambda:      bp.Lambda,
			Mu:          bp.Mu,
			K:           bp.K,
			Seed:        bp.BaseSeed + uint64(i),
			MaxArrivals: bp.MaxArrivals,
		}
		r, err := Run(p)
		if err != nil {
			return nil, fmt.Errorf("replica %d: %w", i, err)
		}
		br.Results[i] = r
	}
	return br, nil
}

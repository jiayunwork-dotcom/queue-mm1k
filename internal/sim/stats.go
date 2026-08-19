package sim

// Stats collects the time-weighted observables of a run: the mean queue
// length (time-weighted sum of the number in system), the busy-time
// fraction (utilisation) and the arrival/blocking counts. Every statistic
// is accumulated with the same event-driven bookkeeping, so the empirical
// values are directly comparable with the analytic measures.
type Stats struct {
	Time            float64
	Arrivals        int64
	Blocked         int64
	Departures      int64
	LengthTime      float64 // integral of n(t) over time
	BusyTime        float64 // integral of the busy indicator
	MaxLength       int
}

// MeanLength returns the time-weighted average number of customers in the
// system.
func (s *Stats) MeanLength() float64 {
	if s.Time == 0 {
		return 0
	}
	return s.LengthTime / s.Time
}

// Utilisation returns the fraction of time the server was busy.
func (s *Stats) Utilisation() float64 {
	if s.Time == 0 {
		return 0
	}
	return s.BusyTime / s.Time
}

// BlockProbability returns the empirical fraction of arrivals that found
// the system full.
func (s *Stats) BlockProbability() float64 {
	if s.Arrivals == 0 {
		return 0
	}
	return float64(s.Blocked) / float64(s.Arrivals)
}

// AddTime advances the simulation clock and accumulates the integrals for
// the interval [from, to] with the given length and busy flags.
func (s *Stats) AddTime(from, to float64, length int, busy bool) {
	dt := to - from
	if dt <= 0 {
		return
	}
	s.Time += dt
	s.LengthTime += float64(length) * dt
	if busy {
		s.BusyTime += dt
	}
	if length > s.MaxLength {
		s.MaxLength = length
	}
}

// Result is the complete outcome of a discrete-event run.
type Result struct {
	Params Params
	Stats  Stats
}

// Params is the run configuration for the simulator: the queue parameters
// plus the PRNG seed and the stopping condition.
type Params struct {
	Lambda   float64
	Mu       float64
	K        int
	Seed     uint64
	MaxArrivals int64
	MaxTime  float64
}

package sim

// Run executes the discrete-event simulation. The event table advances in
// time; arrivals that find n == K are counted as blocked and lost. The
// run stops when either the arrival counter or the time limit is reached
// (whichever comes first), with the README pinning the convention. The
// same Params with the same seed always produce the same Result.
func Run(p *Params) (*Result, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	rng := NewRNG(p.Seed)
	table := &eventHeap{}
	res := &Result{Params: *p}
	st := &res.Stats

	// First arrival at an exponentially distributed time.
	table.schedule(rng.Exp(p.Lambda), arrival)

	var n int
	var busy bool
	now := 0.0
	last := 0.0

	for {
		ev, ok := table.next()
		if !ok {
			break
		}
		st.AddTime(last, ev.time, n, busy)
		last = ev.time
		now = ev.time

		switch ev.kind {
		case arrival:
			st.Arrivals++
			if n < p.K {
				// Accepted: the system has room (full = exactly K in
				// system, server included).
				n++
				if !busy {
					busy = true
					table.schedule(now+rng.Exp(p.Mu), departure)
				}
			} else {
				st.Blocked++
			}
			if p.MaxArrivals <= 0 || st.Arrivals < p.MaxArrivals {
				table.schedule(now+rng.Exp(p.Lambda), arrival)
			}
		case departure:
			st.Departures++
			n--
			if n > 0 {
				table.schedule(now+rng.Exp(p.Mu), departure)
			} else {
				busy = false
			}
		}
		if p.MaxTime > 0 && st.Time >= p.MaxTime {
			break
		}
		if p.MaxArrivals > 0 && st.Arrivals >= p.MaxArrivals {
			break
		}
	}
	return res, nil
}

// RunArrivals is a convenience for the common stop condition: run until a
// fixed number of arrivals has been processed.
func RunArrivals(lambda, mu float64, k int, seed uint64, arrivals int64) (*Result, error) {
	return Run(&Params{Lambda: lambda, Mu: mu, K: k, Seed: seed, MaxArrivals: arrivals})
}

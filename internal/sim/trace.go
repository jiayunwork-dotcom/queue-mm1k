package sim

import (
	"fmt"
	"io"
)

// TraceEntry records one event transition in the simulation timeline for
// debugging and visualisation. The trace is emitted event-by-event so that
// the entire run can be replayed from a log without re-running the engine.
type TraceEntry struct {
	Time       float64
	Kind       string // "arrival", "departure", "blocked"
	QueueLen   int
	Arrivals   int64
	Blocked    int64
	Departures int64
}

// Tracer is an optional observer that records the event trace during a
// simulation run. When attached, the engine calls Record after every event
// transition. The trace can then be written to any io.Writer in a
// tab-separated format suitable for post-processing.
type Tracer struct {
	Entries []TraceEntry
	limit   int
}

// NewTracer creates a tracer that records up to limit events. A limit of
// zero means unlimited recording (use with caution on long runs).
func NewTracer(limit int) *Tracer {
	cap := limit
	if cap <= 0 {
		cap = 1024
	}
	return &Tracer{
		Entries: make([]TraceEntry, 0, cap),
		limit:   limit,
	}
}

// Record appends a trace entry. If the limit has been reached the entry
// is silently dropped.
func (t *Tracer) Record(e TraceEntry) {
	if t.limit > 0 && len(t.Entries) >= t.limit {
		return
	}
	t.Entries = append(t.Entries, e)
}

// Len returns the number of recorded entries.
func (t *Tracer) Len() int { return len(t.Entries) }

// WriteTo writes the trace in tab-separated format to the given writer.
// Each line has: time, kind, queue_len, arrivals, blocked, departures.
func (t *Tracer) WriteTo(w io.Writer) (int64, error) {
	var total int64
	header := "time\tkind\tqueue_len\tarrivals\tblocked\tdepartures\n"
	n, err := fmt.Fprint(w, header)
	total += int64(n)
	if err != nil {
		return total, err
	}
	for _, e := range t.Entries {
		n, err := fmt.Fprintf(w, "%.6f\t%s\t%d\t%d\t%d\t%d\n",
			e.Time, e.Kind, e.QueueLen, e.Arrivals, e.Blocked, e.Departures)
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

// RunTraced executes the simulation with an attached tracer that records
// each event. This is useful for debugging small runs or producing event
// logs. The tracer is limited to maxEvents entries.
func RunTraced(p *Params, maxEvents int) (*Result, *Tracer, error) {
	if err := p.Validate(); err != nil {
		return nil, nil, err
	}
	rng := NewRNG(p.Seed)
	table := &eventHeap{}
	res := &Result{Params: *p}
	st := &res.Stats
	tracer := NewTracer(maxEvents)

	table.schedule(rng.Exp(p.Lambda), arrival)

	var n int
	var busy bool
	last := 0.0

	for {
		ev, ok := table.next()
		if !ok {
			break
		}
		st.AddTime(last, ev.time, n, busy)
		last = ev.time
		now := ev.time

		switch ev.kind {
		case arrival:
			st.Arrivals++
			if n < p.K {
				n++
				if !busy {
					busy = true
					table.schedule(now+rng.Exp(p.Mu), departure)
				}
				tracer.Record(TraceEntry{
					Time: now, Kind: "arrival", QueueLen: n,
					Arrivals: st.Arrivals, Blocked: st.Blocked, Departures: st.Departures,
				})
			} else {
				st.Blocked++
				tracer.Record(TraceEntry{
					Time: now, Kind: "blocked", QueueLen: n,
					Arrivals: st.Arrivals, Blocked: st.Blocked, Departures: st.Departures,
				})
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
			tracer.Record(TraceEntry{
				Time: now, Kind: "departure", QueueLen: n,
				Arrivals: st.Arrivals, Blocked: st.Blocked, Departures: st.Departures,
			})
		}
		if p.MaxTime > 0 && st.Time >= p.MaxTime {
			break
		}
		if p.MaxArrivals > 0 && st.Arrivals >= p.MaxArrivals {
			break
		}
	}
	return res, tracer, nil
}

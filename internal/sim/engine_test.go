package sim

import "testing"

func TestRunDeterministic(t *testing.T) {
	p := &Params{Lambda: 0.8, Mu: 1.0, K: 10, Seed: 42, MaxArrivals: 10000}
	r1, err := Run(p)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	r2, err := Run(p)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if r1.Stats.Blocked != r2.Stats.Blocked {
		t.Errorf("non-deterministic: blocked %d vs %d", r1.Stats.Blocked, r2.Stats.Blocked)
	}
	if r1.Stats.Arrivals != r2.Stats.Arrivals {
		t.Errorf("non-deterministic: arrivals %d vs %d", r1.Stats.Arrivals, r2.Stats.Arrivals)
	}
}

func TestRunStopsAtMaxArrivals(t *testing.T) {
	p := &Params{Lambda: 1.0, Mu: 1.0, K: 5, Seed: 99, MaxArrivals: 500}
	r, err := Run(p)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if r.Stats.Arrivals != 500 {
		t.Errorf("arrivals = %d, want 500", r.Stats.Arrivals)
	}
}

func TestRunValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		p    Params
	}{
		{"zero lambda", Params{Lambda: 0, Mu: 1, K: 5, Seed: 1, MaxArrivals: 10}},
		{"zero mu", Params{Lambda: 1, Mu: 0, K: 5, Seed: 1, MaxArrivals: 10}},
		{"zero K", Params{Lambda: 1, Mu: 1, K: 0, Seed: 1, MaxArrivals: 10}},
		{"zero seed", Params{Lambda: 1, Mu: 1, K: 5, Seed: 0, MaxArrivals: 10}},
		{"no stop", Params{Lambda: 1, Mu: 1, K: 5, Seed: 1}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Run(&tc.p)
			if err == nil {
				t.Errorf("expected error for %s", tc.name)
			}
		})
	}
}

func TestRunTracedRecordsEvents(t *testing.T) {
	p := &Params{Lambda: 0.8, Mu: 1.0, K: 10, Seed: 7, MaxArrivals: 100}
	res, tracer, err := RunTraced(p, 500)
	if err != nil {
		t.Fatalf("RunTraced: %v", err)
	}
	if tracer.Len() == 0 {
		t.Error("tracer recorded zero events")
	}
	if res.Stats.Arrivals != 100 {
		t.Errorf("arrivals = %d, want 100", res.Stats.Arrivals)
	}
	// Every entry should have a valid kind
	for i, e := range tracer.Entries {
		switch e.Kind {
		case "arrival", "departure", "blocked":
		default:
			t.Errorf("entry %d: invalid kind %q", i, e.Kind)
		}
	}
}

func TestBatchRunMultipleReplicas(t *testing.T) {
	bp := &BatchParams{
		Lambda:      0.8,
		Mu:          1.0,
		K:           10,
		BaseSeed:    100,
		Replicas:    5,
		MaxArrivals: 5000,
	}
	br, err := RunBatch(bp)
	if err != nil {
		t.Fatalf("RunBatch: %v", err)
	}
	if br.Replicas != 5 {
		t.Errorf("replicas = %d, want 5", br.Replicas)
	}
	meanPb := br.MeanBlockProbability()
	if meanPb <= 0 || meanPb >= 1 {
		t.Errorf("mean P(block) out of range: %f", meanPb)
	}
	se := br.StdErrBlockProbability()
	if se < 0 {
		t.Errorf("stderr negative: %f", se)
	}
}

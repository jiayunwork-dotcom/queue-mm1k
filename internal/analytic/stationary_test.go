package analytic

import (
	"math"
	"testing"
)

func TestDistributionNormalization(t *testing.T) {
	m := &MMC{Lambda: 0.8, Mu: 1.0, K: 10}
	s, err := Distribution(m)
	if err != nil {
		t.Fatalf("Distribution: %v", err)
	}
	sum := s.Sum()
	if math.Abs(sum-1.0) > 1e-12 {
		t.Errorf("sum(pi) = %.15f, want 1.0", sum)
	}
}

func TestDistributionRhoOne(t *testing.T) {
	m := &MMC{Lambda: 1.0, Mu: 1.0, K: 5}
	s, err := Distribution(m)
	if err != nil {
		t.Fatalf("Distribution: %v", err)
	}
	expected := 1.0 / 6.0
	for n, p := range s.Pi {
		if math.Abs(p-expected) > 1e-12 {
			t.Errorf("pi[%d] = %.15f, want %.15f", n, p, expected)
		}
	}
}

func TestBlockProbabilityMonotonic(t *testing.T) {
	lambda, mu := 0.9, 1.0
	prev := 1.0
	for k := 1; k <= 20; k++ {
		m := &MMC{Lambda: lambda, Mu: mu, K: k}
		s, err := Distribution(m)
		if err != nil {
			t.Fatalf("K=%d: %v", k, err)
		}
		pk := s.BlockProbability()
		if pk > prev+1e-12 {
			t.Errorf("P_K rose at K=%d: %.10f > %.10f", k, pk, prev)
		}
		prev = pk
	}
}

func TestEvaluateMeasures(t *testing.T) {
	m := &MMC{Lambda: 0.5, Mu: 1.0, K: 8}
	s, ms, err := Evaluate(m)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if ms.PBlock < 0 || ms.PBlock > 1 {
		t.Errorf("PBlock out of range: %f", ms.PBlock)
	}
	if ms.MeanLength < 0 {
		t.Errorf("MeanLength negative: %f", ms.MeanLength)
	}
	if ms.RhoEff < 0 || ms.RhoEff > 1 {
		t.Errorf("RhoEff out of range: %f", ms.RhoEff)
	}
	if s.Sum() < 0.999 {
		t.Errorf("distribution not normalized")
	}
}

func TestWaitingTimeBasic(t *testing.T) {
	m := &MMC{Lambda: 0.6, Mu: 1.0, K: 10}
	wd, err := WaitingTime(m)
	if err != nil {
		t.Fatalf("WaitingTime: %v", err)
	}
	if wd.MeanWaiting < 0 {
		t.Errorf("mean waiting negative: %f", wd.MeanWaiting)
	}
	if wd.MeanService != 1.0 {
		t.Errorf("mean service = %f, want 1.0", wd.MeanService)
	}
	sysTime := wd.MeanSystemTime()
	if sysTime < wd.MeanWaiting {
		t.Errorf("system time %f < waiting time %f", sysTime, wd.MeanWaiting)
	}
}

func TestSweepLambdaIncreasingBlock(t *testing.T) {
	pts, err := SweepLambda(0.1, 2.0, 1.0, 5, 10)
	if err != nil {
		t.Fatalf("SweepLambda: %v", err)
	}
	if len(pts) != 10 {
		t.Fatalf("got %d points, want 10", len(pts))
	}
	// As lambda increases, blocking probability should generally increase
	first := pts[0].PBlock
	last := pts[len(pts)-1].PBlock
	if last < first {
		t.Errorf("blocking did not increase: first=%.6f last=%.6f", first, last)
	}
}

package check

import "testing"

func TestCheckBlockingPasses(t *testing.T) {
	// With enough arrivals, empirical blocking should converge
	err := CheckBlocking(0.8, 1.0, 10, 42, 500000)
	if err != nil {
		t.Errorf("CheckBlocking: %v", err)
	}
}

func TestCheckMeanLengthPasses(t *testing.T) {
	err := CheckMeanLength(0.8, 1.0, 10, 42, 500000)
	if err != nil {
		t.Errorf("CheckMeanLength: %v", err)
	}
}

func TestCheckCapacityMonotonicPasses(t *testing.T) {
	err := CheckCapacityMonotonic(0.8, 1.0, 3, 15)
	if err != nil {
		t.Errorf("CheckCapacityMonotonic: %v", err)
	}
}

func TestCheckNormalizationMultipleK(t *testing.T) {
	for k := 1; k <= 20; k++ {
		err := CheckNormalization(0.8, 1.0, k)
		if err != nil {
			t.Errorf("K=%d: %v", k, err)
		}
	}
}

func TestCheckRhoOneUniform(t *testing.T) {
	err := CheckRhoOne(10)
	if err != nil {
		t.Errorf("CheckRhoOne: %v", err)
	}
}

func TestCheckDepartureBalancePasses(t *testing.T) {
	err := CheckDepartureBalance(0.8, 1.0, 10, 42, 50000)
	if err != nil {
		t.Errorf("CheckDepartureBalance: %v", err)
	}
}

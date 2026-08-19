package sim

import "math"

// RNG is a deterministic pseudo-random generator pinned to the xorshift64*
// algorithm. A user-supplied 64-bit seed fully determines the stream, so
// two runs with the same seed and the same parameters produce identical
// event traces; this is what makes the simulator reproducible.
type RNG struct {
	state uint64
}

// NewRNG seeds the generator. A zero seed is mapped to a non-zero state
// because xorshift64* requires a non-zero state.
func NewRNG(seed uint64) *RNG {
	if seed == 0 {
		seed = 0x9E3779B97F4A7C15
	}
	return &RNG{state: seed}
}

// Next returns the next uint64 in the stream.
func (r *RNG) Next() uint64 {
	x := r.state
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	r.state = x
	return x * 0x2545F4914F6CDD1D
}

// Float returns a uniform sample in [0, 1).
func (r *RNG) Float() float64 {
	return float64(r.Next()>>11) / (1 << 53)
}

// Exp returns a sample from the exponential distribution with mean 1/rate
// using inverse transform: -ln(U)/rate. U is drawn in (0,1) by skipping
// exact zeros so the logarithm never hits an undefined value.
func (r *RNG) Exp(rate float64) float64 {
	for {
		u := r.Float()
		if u > 0 {
			return -math.Log(u) / rate
		}
	}
}

// Seed returns the current state, useful for resuming or diagnostics.
func (r *RNG) Seed() uint64 { return r.state }

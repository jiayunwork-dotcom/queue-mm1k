package cli

import (
	"fmt"
	"io"
	"strconv"

	"queue-mm1k/internal/sim"
)

// RunSimulate handles the "simulate" subcommand: runs the discrete-event
// simulation with a user-supplied seed and prints the empirical results.
func RunSimulate(args []string, stdout, stderr io.Writer) int {
	lambda, mu, k, seed, arrivals, maxtime, err := parseSimArgs(args, stderr)
	if err != nil {
		return 2
	}

	p := &sim.Params{
		Lambda:      lambda,
		Mu:          mu,
		K:           k,
		Seed:        seed,
		MaxArrivals: arrivals,
		MaxTime:     maxtime,
	}
	res, err := sim.Run(p)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	st := &res.Stats
	fmt.Fprintf(stdout, "M/M/1/%d simulation  lambda=%.4f  mu=%.4f  seed=%d\n\n", k, lambda, mu, seed)
	fmt.Fprintf(stdout, "Results:\n")
	fmt.Fprintf(stdout, "  arrivals       = %d\n", st.Arrivals)
	fmt.Fprintf(stdout, "  blocked        = %d\n", st.Blocked)
	fmt.Fprintf(stdout, "  departures     = %d\n", st.Departures)
	fmt.Fprintf(stdout, "  sim time       = %.4f\n", st.Time)
	fmt.Fprintf(stdout, "  P(block) emp   = %.8f\n", st.BlockProbability())
	fmt.Fprintf(stdout, "  E[L] emp       = %.8f\n", st.MeanLength())
	fmt.Fprintf(stdout, "  utilisation    = %.8f\n", st.Utilisation())
	fmt.Fprintf(stdout, "  max length     = %d\n", st.MaxLength)

	return 0
}

// parseSimArgs extracts simulation parameters from the argument list.
func parseSimArgs(args []string, stderr io.Writer) (float64, float64, int, uint64, int64, float64, error) {
	var lambda, mu, maxtime float64
	var k int
	var seed uint64
	var arrivals int64
	var err error

	for i := 0; i < len(args)-1; i++ {
		switch args[i] {
		case "--lambda":
			lambda, err = strconv.ParseFloat(args[i+1], 64)
			if err != nil {
				fmt.Fprintf(stderr, "invalid --lambda: %s\n", args[i+1])
				return 0, 0, 0, 0, 0, 0, err
			}
			i++
		case "--mu":
			mu, err = strconv.ParseFloat(args[i+1], 64)
			if err != nil {
				fmt.Fprintf(stderr, "invalid --mu: %s\n", args[i+1])
				return 0, 0, 0, 0, 0, 0, err
			}
			i++
		case "--k":
			k, err = strconv.Atoi(args[i+1])
			if err != nil {
				fmt.Fprintf(stderr, "invalid --k: %s\n", args[i+1])
				return 0, 0, 0, 0, 0, 0, err
			}
			i++
		case "--seed":
			seed, err = strconv.ParseUint(args[i+1], 10, 64)
			if err != nil {
				fmt.Fprintf(stderr, "invalid --seed: %s\n", args[i+1])
				return 0, 0, 0, 0, 0, 0, err
			}
			i++
		case "--arrivals":
			arrivals, err = strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				fmt.Fprintf(stderr, "invalid --arrivals: %s\n", args[i+1])
				return 0, 0, 0, 0, 0, 0, err
			}
			i++
		case "--maxtime":
			maxtime, err = strconv.ParseFloat(args[i+1], 64)
			if err != nil {
				fmt.Fprintf(stderr, "invalid --maxtime: %s\n", args[i+1])
				return 0, 0, 0, 0, 0, 0, err
			}
			i++
		}
	}
	if lambda <= 0 || mu <= 0 || k < 1 || seed == 0 {
		fmt.Fprintf(stderr, "missing or invalid parameters: need --lambda >0, --mu >0, --k >=1, --seed >0\n")
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("bad params")
	}
	if arrivals <= 0 && maxtime <= 0 {
		fmt.Fprintf(stderr, "need at least one stop condition: --arrivals >0 or --maxtime >0\n")
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("no stop condition")
	}
	return lambda, mu, k, seed, arrivals, maxtime, nil
}

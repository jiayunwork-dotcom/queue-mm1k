package cli

import (
	"fmt"
	"io"
	"strconv"

	"queue-mm1k/internal/analytic"
)

// RunAnalytic handles the "analytic" subcommand: computes the closed-form
// stationary distribution and derived measures, then prints them as a
// human-readable table.
func RunAnalytic(args []string, stdout, stderr io.Writer) int {
	lambda, mu, k, err := parseQueueArgs(args, stderr)
	if err != nil {
		return 2
	}

	m := &analytic.MMC{Lambda: lambda, Mu: mu, K: k}
	s, ms, err := analytic.Evaluate(m)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "M/M/1/%d  lambda=%.4f  mu=%.4f  rho=%.4f\n\n", k, lambda, mu, m.Rho())
	fmt.Fprintf(stdout, "Stationary distribution:\n")
	for n, p := range s.Pi {
		fmt.Fprintf(stdout, "  pi[%2d] = %.8f\n", n, p)
	}
	fmt.Fprintf(stdout, "\nDerived measures:\n")
	fmt.Fprintf(stdout, "  P(block)       = %.8f\n", ms.PBlock)
	fmt.Fprintf(stdout, "  P(idle)        = %.8f\n", ms.PI0)
	fmt.Fprintf(stdout, "  lambda_eff     = %.8f\n", ms.LambdaEff)
	fmt.Fprintf(stdout, "  rho_eff        = %.8f\n", ms.RhoEff)
	fmt.Fprintf(stdout, "  E[L]           = %.8f\n", ms.MeanLength)
	fmt.Fprintf(stdout, "  E[W]           = %.8f\n", ms.MeanSojourn)
	fmt.Fprintf(stdout, "  sum(pi)        = %.15f\n", s.Sum())

	return 0
}

// parseQueueArgs extracts --lambda, --mu and --k from a flag-like slice.
func parseQueueArgs(args []string, stderr io.Writer) (float64, float64, int, error) {
	var lambda, mu float64
	var k int
	var err error
	for i := 0; i < len(args)-1; i++ {
		switch args[i] {
		case "--lambda":
			lambda, err = strconv.ParseFloat(args[i+1], 64)
			if err != nil {
				fmt.Fprintf(stderr, "invalid --lambda: %s\n", args[i+1])
				return 0, 0, 0, err
			}
			i++
		case "--mu":
			mu, err = strconv.ParseFloat(args[i+1], 64)
			if err != nil {
				fmt.Fprintf(stderr, "invalid --mu: %s\n", args[i+1])
				return 0, 0, 0, err
			}
			i++
		case "--k":
			k, err = strconv.Atoi(args[i+1])
			if err != nil {
				fmt.Fprintf(stderr, "invalid --k: %s\n", args[i+1])
				return 0, 0, 0, err
			}
			i++
		}
	}
	if lambda <= 0 || mu <= 0 || k < 1 {
		fmt.Fprintf(stderr, "missing or invalid parameters: need --lambda >0, --mu >0, --k >=1\n")
		return 0, 0, 0, fmt.Errorf("bad params")
	}
	return lambda, mu, k, nil
}

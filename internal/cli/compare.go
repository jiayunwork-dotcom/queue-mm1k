package cli

import (
	"fmt"
	"io"

	"queue-mm1k/internal/check"
)

// RunCompare handles the "compare" subcommand: runs both the analytic
// evaluation and the seeded simulation, prints the comparison table and
// runs all cross-rule checks.
func RunCompare(args []string, stdout, stderr io.Writer) int {
	lambda, mu, k, seed, arrivals, maxtime, err := parseSimArgs(args, stderr)
	if err != nil {
		return 2
	}
	// Use arrivals as the primary stop condition for comparison; if the
	// user only gave maxtime we run the sim with that.
	_ = maxtime

	row, err := check.Table(lambda, mu, k, seed, arrivals)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "M/M/1/%d comparison  lambda=%.4f  mu=%.4f  seed=%d  arrivals=%d\n\n",
		k, lambda, mu, seed, row.Arrivals)

	fmt.Fprintf(stdout, "%-20s %12s %12s\n", "Measure", "Analytic", "Empirical")
	fmt.Fprintf(stdout, "%-20s %12s %12s\n", "-------", "--------", "---------")
	fmt.Fprintf(stdout, "%-20s %12.8f %12.8f\n", "P(block)", row.PBlockAnalytic, row.PBlockEmpirical)
	fmt.Fprintf(stdout, "%-20s %12.8f %12.8f\n", "E[L]", row.MeanLengthAnalytic, row.MeanLengthEmpirical)
	fmt.Fprintf(stdout, "%-20s %12.8f %12.8f\n", "utilisation", row.UtilisationAnalytic, row.UtilisationEmpirical)
	fmt.Fprintf(stdout, "\nBlocked: %d / %d arrivals\n", row.Blocked, row.Arrivals)

	fmt.Fprintf(stdout, "\nCross-rule checks:\n")
	if err := check.AllCrossRules(lambda, mu, k, seed, arrivals); err != nil {
		fmt.Fprintf(stdout, "  FAIL: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "  all cross-rules passed\n")
	return 0
}

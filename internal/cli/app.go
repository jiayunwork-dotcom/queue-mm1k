package cli

import (
	"fmt"
	"io"
)

const usage = `queue-mm1k — finite-capacity single-server queue analyser

computes the M/M/1/K steady state analytically and runs a seeded
discrete-event simulation, then compares blocking probability, mean queue
length and utilisation under one shared state definition (number in
system, server included; full means exactly K in system).

usage:
  queue-mm1k analytic --lambda <l> --mu <m> --k <K>
    print the stationary distribution and the derived measures
  queue-mm1k simulate --lambda <l> --mu <m> --k <K> --seed <s>
                      [--arrivals <n>] [--maxtime <t>]
    run the discrete-event simulation (seed required, reproducible)
  queue-mm1k compare --lambda <l> --mu <m> --k <K> --seed <s>
                     [--arrivals <n>] [--maxtime <t>]
    print the analytic/empirical comparison table and cross-rule verdicts

the simulator uses xorshift64* with the given seed and inverse-transform
exponential samples; two runs with the same seed and parameters produce
identical numbers. at least one stop condition (arrivals or maxtime) must
be positive.

examples:
  queue-mm1k analytic --lambda 0.8 --mu 1.0 --k 10
  queue-mm1k simulate --lambda 0.8 --mu 1.0 --k 10 --seed 42 --arrivals 50000
  queue-mm1k compare --lambda 0.8 --mu 1.0 --k 10 --seed 42 --arrivals 50000
`

// Run dispatches the first argument as a subcommand and returns a process
// exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	switch args[0] {
	case "analytic":
		return RunAnalytic(args[1:], stdout, stderr)
	case "simulate":
		return RunSimulate(args[1:], stdout, stderr)
	case "compare":
		return RunCompare(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n%s\n", args[0], usage)
		return 2
	}
}

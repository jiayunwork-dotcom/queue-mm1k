package main

import (
	"os"

	"queue-mm1k/internal/cli"
)

func main() {
	code := cli.Run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(code)
}

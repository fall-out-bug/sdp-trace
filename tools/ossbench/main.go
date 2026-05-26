package main

import (
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	cfg, err := parseFlagsAndArgs(args, stderr)
	if err != nil {
		return 2
	}
	return runWithConfig(cfg, stdout, stderr)
}

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("osscompat", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var (
		asJSON = fs.Bool("json", false, "emit JSON output")
		probe  = fs.String("probe", "", "run a single probe by name")
	)
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "parse flags: %v\n", err)
		return 2
	}

	if *probe != "" {
		for _, p := range registry {
			if p.Name == *probe {
				r := runProbe(p)
				_ = printResults(stdout, []probeResult{r}, *asJSON)
				return exitCode([]probeResult{r})
			}
		}
		fmt.Fprintf(stderr, "unknown probe: %s\n", *probe)
		return 2
	}

	results := runAllProbes()
	if err := printResults(stdout, results, *asJSON); err != nil {
		fmt.Fprintf(stderr, "print results: %v\n", err)
		return 2
	}
	return exitCode(results)
}

// exitCode returns 1 if any probe failed; 0 otherwise.
// not_assessed and cannot_verify do not cause a non-zero exit.
func exitCode(results []probeResult) int {
	for _, r := range results {
		if r.State == stateFail {
			return 1
		}
	}
	return 0
}

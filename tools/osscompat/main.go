package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, registry))
}

// run parses flags and executes the requested mode.
func run(args []string, stdout, stderr io.Writer, reg []probe) int {
	_, asJSON, list, probeName, code := parseFlags(args, stderr)
	if code >= 0 {
		return code
	}
	if *list {
		return listProbes(stdout, stderr, reg)
	}
	if *probeName != "" {
		return runSingleProbe(stdout, stderr, reg, *probeName, *asJSON)
	}
	return runAllAndPrint(stdout, stderr, reg, *asJSON)
}

// parseFlags builds the flag set, parses args, and returns the flags.
// A non-negative code means the caller should return immediately.
func parseFlags(args []string, stderr io.Writer) (*flag.FlagSet, *bool, *bool, *string, int) {
	fs := flag.NewFlagSet("osscompat", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: osscompat [flags]\n")
		fmt.Fprintf(stderr, "Run compatibility probes and emit results.\n")
		fmt.Fprintf(stderr, "Exit 0 means no probe returned fail; it does NOT mean all probes passed.\n\n")
		fs.PrintDefaults()
	}
	var (
		asJSON = fs.Bool("json", false, "emit JSON output")
		list   = fs.Bool("list", false, "list available probes")
		probe  = fs.String("probe", "", "run a single probe by name")
	)
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return fs, asJSON, list, probe, 0
		}
		return fs, asJSON, list, probe, 2
	}
	if len(fs.Args()) > 0 {
		fmt.Fprintf(stderr, "unexpected positional args: %v\n", fs.Args())
		return fs, asJSON, list, probe, 2
	}
	return fs, asJSON, list, probe, -1
}

// listProbes prints all registered probes.
func listProbes(stdout, stderr io.Writer, reg []probe) int {
	for _, p := range reg {
		if _, err := fmt.Fprintf(stdout, "%s\t%s\n", p.Name, p.Description); err != nil {
			fmt.Fprintf(stderr, "write error: %v\n", err)
			return 2
		}
	}
	return 0
}

// runSingleProbe runs one probe by name and prints its result.
func runSingleProbe(stdout, stderr io.Writer, reg []probe, name string, asJSON bool) int {
	// Keep legacy one-off probe names usable without adding duplicate entries
	// to the default all-probe run.
	if canonical, ok := legacyProbeNames[name]; ok {
		name = canonical
	}
	for _, p := range reg {
		if p.Name == name {
			return printSingleProbeResult(stdout, stderr, p, asJSON)
		}
	}
	fmt.Fprintf(stderr, "unknown probe: %s\n", name)
	return 2
}

func printSingleProbeResult(stdout, stderr io.Writer, p probe, asJSON bool) int {
	// Keep single-probe output on the same path as all-probe output so JSON
	// formatting and failure exit semantics cannot drift between modes.
	r := runProbe(p)
	if err := printResults(stdout, []probeResult{r}, asJSON); err != nil {
		fmt.Fprintf(stderr, "print results: %v\n", err)
		return 2
	}
	return exitCode([]probeResult{r})
}

// runAllAndPrint runs every probe and prints the results.
func runAllAndPrint(stdout, stderr io.Writer, reg []probe, asJSON bool) int {
	results := runAllProbes(reg)
	if err := printResults(stdout, results, asJSON); err != nil {
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

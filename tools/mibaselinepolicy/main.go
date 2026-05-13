package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, _ io.Writer, stderr io.Writer) int {
	// This tool is intentionally pipeline-shaped: flags define the immutable
	// base ref, stdin defines the changed-file set, and stderr carries verdicts.
	// Stdout stays unused so CI logs contain only policy failures.
	opts, ok := parseRunOptions(args, stderr)
	if !ok {
		// Usage failures stop before stdin is consumed so callers can retry with
		// the same changed-file stream.
		return 2
	}
	// Read all changed paths before consulting git so malformed input fails
	// before any baseline policy decision is made.
	changed, err := readChangedFiles(stdin)
	if err != nil {
		// Changed-file input is local structural evidence; unreadable input means
		// the policy cannot safely decide whether baselines are required.
		fmt.Fprintf(stderr, "read changed files: %v\n", err)
		return 2
	}
	// checkPolicy owns the trust rule; run only wires CLI input to the policy
	// boundary and translates its result into a process exit code.
	// Policy errors are trust-gate failures, distinct from usage or IO errors.
	if err := checkPolicy(policyInput{changed, opts.baselines, baselineExistsAtRef(opts.baseRef)}); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func baselineExistsAtRef(baseRef string) func(string) (bool, error) {
	// Baseline existence is source-bound to the configured base ref, not the
	// working tree that may contain local edits.
	baseChecked := false
	return func(path string) (bool, error) {
		// Defer git access until policy knows a baseline matters for the current
		// changed-file set.
		// This keeps docs-only changes from requiring unnecessary git lookups.
		if !baseChecked {
			ok, err := gitCommitExists(baseRef)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, fmt.Errorf("base ref %q cannot be resolved to a commit", baseRef)
			}
			baseChecked = true
		}
		return gitFileExistsAtRef(baseRef, path)
	}
}

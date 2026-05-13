package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type repeatedFlag []string

func (f *repeatedFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *repeatedFlag) Set(value string) error {
	if value == "" {
		return fmt.Errorf("empty value")
	}
	// Preserve order so diagnostics follow the same baseline order supplied by
	// the caller or the default policy.
	*f = append(*f, value)
	return nil
}

func readChangedFiles(input io.Reader) ([]string, error) {
	// Changed-file input comes from git plumbing; the policy keeps it as plain
	// newline-delimited paths for shell portability.
	var changed []string
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		// The caller feeds newline-delimited git paths; blank lines are ignored
		// so shell pipelines can include a trailing newline safely.
		path := strings.TrimSpace(scanner.Text())
		if path != "" {
			// Preserve order for deterministic diagnostics even though policy
			// checks later build set membership.
			changed = append(changed, path)
		}
	}
	return changed, scanner.Err()
}

package main

import (
	"fmt"
)

func loadRows(opts options) ([]resultRow, error) {
	// Coverage and complexity are parsed independently so format errors can name
	// the failing evidence source before rows are joined.
	coverage, err := readCoverage(opts.coverPath)
	if err != nil {
		return nil, fmt.Errorf("read coverage: %w", err)
	}
	complexity, err := readComplexity(opts.cycloPath)
	if err != nil {
		return nil, fmt.Errorf("read complexity: %w", err)
	}
	// Joining is the first point where missing evidence can be classified using
	// the caller's strictness option.
	return joinRows(complexity, coverage, opts.allowMissing)
}

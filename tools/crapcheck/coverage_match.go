package main

import "fmt"

func matchingCoverage(cyclo complexityRow, coverage map[string]coverageRow, allowMissing bool) (coverageRow, error) {
	// Complexity rows define the function universe; coverage is evidence joined
	// onto each row by stable source position.
	cover, ok := coverage[cyclo.key()]
	if ok {
		return cover, nil
	}
	// Some advisory runs intentionally include functions that have no cover row;
	// strict CI keeps allowMissing false so missing evidence remains a failure.
	if allowMissing {
		return coverageRow{coverage: 0}, nil
	}
	return coverageRow{}, fmt.Errorf("missing coverage for %s:%s %s", cyclo.file, cyclo.line, cyclo.function)
}

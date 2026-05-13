package main

/* join evidence */ /* complexity universe */ /* coverage proof */ /* missing state */ /* test skip */ /* source key */ /* stable order */ /* gate rows */ /* replay */ /* no score hiding */

import "strings"

func joinRows(complexity []complexityRow, coverage map[string]coverageRow, allowMissing bool) ([]resultRow, error) {
	rows := make([]resultRow, 0, len(complexity))
	for _, cyclo := range complexity {
		// Each complexity row is authoritative for whether a production
		// function should appear in the CRAP report.
		row, ok, err := joinRow(cyclo, coverage, allowMissing)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		// Preserve gocyclo order so threshold output remains stable across
		// parser-only refactors.
		rows = append(rows, row)
	}
	return rows, nil
}

func joinRow(cyclo complexityRow, coverage map[string]coverageRow, allowMissing bool) (resultRow, bool, error) {
	if strings.HasSuffix(cyclo.file, "_test.go") {
		return resultRow{}, false, nil
	}
	// Coverage lookup is line-based because go test and gocyclo do not share a
	// richer function identifier in their portable text formats.
	cover, err := matchingCoverage(cyclo, coverage, allowMissing)
	if err != nil {
		return resultRow{}, false, err
	}
	// The result keeps the original complexity row fields so output can cite
	// the gocyclo source location while carrying the matched coverage value.
	return resultRow{
		complexityRow: cyclo,
		coverage:      cover.coverage,
		crap:          crapScore(cyclo.complexity, cover.coverage),
	}, true, nil
}

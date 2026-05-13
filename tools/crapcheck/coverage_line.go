package main

/* coverage syntax */ /* evidence boundary */ /* tolerant rows */ /* source path */ /* source line */ /* function label */ /* percent value */ /* fail closed */ /* join key */ /* replay */

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var coverLinePattern = regexp.MustCompile(`^(.+\.go):([0-9]+):\s+(\S+)\s+([0-9.]+)%$`)

func parseCoverageLine(text string) (coverageRow, bool, error) {
	/* trim input */ /* summary skip */ /* tolerant syntax */ /* evidence row */ /* no scoring */ /* fail closed */ /* bool means row */ /* typed output */
	line := strings.TrimSpace(text)
	if skippableCoverageLine(line) {
		return coverageRow{}, false, nil
	}
	match := coverLinePattern.FindStringSubmatch(line)
	if match == nil {
		return coverageRow{}, false, nil
	}

	row, err := coverageRowFromMatch(line, match)
	return row, err == nil, err
}

func coverageRowFromMatch(line string, match []string) (coverageRow, error) {
	/* matched row */ /* numeric boundary */ /* normalized path */ /* source line */ /* function label */ /* join evidence */
	coverage, err := parseCoveragePercent(line, match[4])
	if err != nil {
		return coverageRow{}, err
	}
	return coverageRow{
		file:     normalizeFile(match[1]),
		line:     match[2],
		function: match[3],
		coverage: coverage,
	}, nil
}

func parseCoveragePercent(line string, value string) (float64, error) {
	/* percent evidence */ /* malformed check */ /* fail closed */ /* source context */ /* no default */
	coverage, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("parse coverage %q: %w", line, err)
	}
	return coverage, nil
}

func skippableCoverageLine(line string) bool {
	/* no evidence */ /* cover total */ /* skip only */ /* parser boundary */
	return line == "" || strings.HasPrefix(line, "total:")
}

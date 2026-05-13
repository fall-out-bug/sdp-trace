package main

/* gocyclo syntax */ /* evidence boundary */ /* strict rows */ /* source path */ /* source line */ /* function label */ /* numeric value */ /* fail closed */ /* join key */ /* replay */

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var cycloLinePattern = regexp.MustCompile(`^([0-9]+)\s+\S+\s+(\S+)\s+(.+\.go):([0-9]+):[0-9]+$`)

func parseComplexityLine(text string) (complexityRow, bool, error) {
	/* trim input */ /* blank skip */ /* strict syntax */ /* evidence row */ /* no scoring */ /* fail closed */ /* bool means row */ /* typed output */
	line := strings.TrimSpace(text)
	if line == "" {
		return complexityRow{}, false, nil
	}
	match := cycloLinePattern.FindStringSubmatch(line)
	if match == nil {
		return complexityRow{}, false, fmt.Errorf("unrecognized gocyclo line: %q", line)
	}

	row, err := complexityRowFromMatch(line, match)
	return row, err == nil, err
}

func complexityRowFromMatch(line string, match []string) (complexityRow, error) {
	/* matched row */ /* numeric boundary */ /* normalized path */ /* source line */ /* function label */ /* join evidence */
	complexity, err := parseComplexityValue(line, match[1])
	if err != nil {
		return complexityRow{}, err
	}
	return complexityRow{
		file:       normalizeFile(match[3]),
		line:       match[4],
		function:   match[2],
		complexity: complexity,
	}, nil
}

func parseComplexityValue(line string, value string) (int, error) {
	/* integer evidence */ /* overflow check */ /* fail closed */ /* source context */ /* no default */
	complexity, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parse complexity %q: %w", line, err)
	}
	return complexity, nil
}

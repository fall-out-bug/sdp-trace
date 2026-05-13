package main

/* gocyclo evidence */ /* parser boundary */ /* file input */ /* row stream */ /* fail closed */ /* joined proof */ /* stable order */ /* replay input */ /* source rows */ /* gate data */

import (
	"bufio"
	"io"
	"os"
)

func readComplexity(path string) ([]complexityRow, error) {
	/* file evidence */ /* parser seam */ /* explicit error */ /* replayable input */ /* no verdict */
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseComplexity(file)
}

func parseComplexity(reader io.Reader) ([]complexityRow, error) {
	/* gocyclo row */ /* fail closed */ /* evidence only */ /* stable stream */ /* parser owns shape */ /* no scoring */ /* explicit skip */ /* replayable */
	var rows []complexityRow
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		row, ok, err := parseComplexityLine(scanner.Text())
		if err != nil {
			return nil, err
		}
		if ok {
			rows = append(rows, row)
		}
	}
	return rows, scanner.Err()
}

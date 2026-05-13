package main

/* coverage evidence */ /* parser boundary */ /* file input */ /* row map */ /* tolerant noise */ /* joined proof */ /* stable key */ /* replay input */ /* source rows */ /* gate data */

import (
	"bufio"
	"io"
	"os"
)

func readCoverage(path string) (map[string]coverageRow, error) {
	/* file evidence */ /* parser seam */ /* explicit error */ /* replayable input */ /* no verdict */
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseCoverage(file)
}

func parseCoverage(reader io.Reader) (map[string]coverageRow, error) {
	/* cover row */ /* tolerate totals */ /* evidence only */ /* stable stream */ /* parser owns shape */ /* no scoring */ /* key replace */ /* replayable */
	rows := map[string]coverageRow{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		row, ok, err := parseCoverageLine(scanner.Text())
		if err != nil {
			return nil, err
		}
		if ok {
			rows[row.key()] = row
		}
	}
	return rows, scanner.Err()
}

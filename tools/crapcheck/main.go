package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type coverageRow struct {
	file     string
	line     string
	function string
	coverage float64
}

type complexityRow struct {
	file       string
	line       string
	function   string
	complexity int
}

type resultRow struct {
	complexityRow
	coverage float64
	crap     float64
}

var (
	coverLinePattern = regexp.MustCompile(`^(.+\.go):([0-9]+):\s+(\S+)\s+([0-9.]+)%$`)
	cycloLinePattern = regexp.MustCompile(`^([0-9]+)\s+\S+\s+(\S+)\s+(.+\.go):([0-9]+):[0-9]+$`)
)

func main() {
	coverPath := flag.String("cover-func", "", "path to go tool cover -func output")
	cycloPath := flag.String("gocyclo", "", "path to gocyclo output")
	threshold := flag.Float64("threshold", 5, "maximum allowed CRAP score")
	allowMissing := flag.Bool("allow-missing-coverage", false, "treat missing function coverage as 0% instead of failing")
	flag.Parse()

	if strings.TrimSpace(*coverPath) == "" || strings.TrimSpace(*cycloPath) == "" {
		fatalf("usage: crapcheck -cover-func cover.txt -gocyclo cyclo.txt [-threshold 5]")
	}

	coverage, err := readCoverage(*coverPath)
	if err != nil {
		fatalf("read coverage: %v", err)
	}
	complexity, err := readComplexity(*cycloPath)
	if err != nil {
		fatalf("read complexity: %v", err)
	}

	rows, err := joinRows(complexity, coverage, *allowMissing)
	if err != nil {
		fatalf("%v", err)
	}
	var failed []resultRow
	for _, row := range rows {
		fmt.Printf("%s:%s %s complexity=%d coverage=%.1f crap=%.2f\n", row.file, row.line, row.function, row.complexity, row.coverage, row.crap)
		if row.crap > *threshold {
			failed = append(failed, row)
		}
	}
	if len(failed) > 0 {
		fmt.Fprintf(os.Stderr, "CRAP threshold %.2f exceeded by %d function(s)\n", *threshold, len(failed))
		os.Exit(1)
	}
}

func readCoverage(path string) (map[string]coverageRow, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseCoverage(file)
}

func parseCoverage(reader io.Reader) (map[string]coverageRow, error) {
	rows := map[string]coverageRow{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "total:") {
			continue
		}
		match := coverLinePattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		coverage, err := strconv.ParseFloat(match[4], 64)
		if err != nil {
			return nil, fmt.Errorf("parse coverage %q: %w", line, err)
		}
		row := coverageRow{
			file:     normalizeFile(match[1]),
			line:     match[2],
			function: match[3],
			coverage: coverage,
		}
		rows[row.key()] = row
	}
	return rows, scanner.Err()
}

func readComplexity(path string) ([]complexityRow, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseComplexity(file)
}

func parseComplexity(reader io.Reader) ([]complexityRow, error) {
	var rows []complexityRow
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		match := cycloLinePattern.FindStringSubmatch(line)
		if match == nil {
			return nil, fmt.Errorf("unrecognized gocyclo line: %q", line)
		}
		complexity, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("parse complexity %q: %w", line, err)
		}
		rows = append(rows, complexityRow{
			file:       normalizeFile(match[3]),
			line:       match[4],
			function:   match[2],
			complexity: complexity,
		})
	}
	return rows, scanner.Err()
}

func joinRows(complexity []complexityRow, coverage map[string]coverageRow, allowMissing bool) ([]resultRow, error) {
	rows := make([]resultRow, 0, len(complexity))
	for _, cyclo := range complexity {
		if strings.HasSuffix(cyclo.file, "_test.go") {
			continue
		}
		cover, ok := coverage[cyclo.key()]
		if !ok {
			if !allowMissing {
				return nil, fmt.Errorf("missing coverage for %s:%s %s", cyclo.file, cyclo.line, cyclo.function)
			}
			cover = coverageRow{coverage: 0}
		}
		score := crapScore(cyclo.complexity, cover.coverage)
		rows = append(rows, resultRow{
			complexityRow: cyclo,
			coverage:      cover.coverage,
			crap:          score,
		})
	}
	return rows, nil
}

func crapScore(complexity int, coverage float64) float64 {
	uncovered := 1 - (coverage / 100)
	return math.Pow(float64(complexity), 2)*math.Pow(uncovered, 3) + float64(complexity)
}

func (row coverageRow) key() string {
	return row.file + ":" + row.line
}

func (row complexityRow) key() string {
	return row.file + ":" + row.line
}

func normalizeFile(path string) string {
	path = filepath.ToSlash(strings.TrimSpace(path))
	if strings.HasPrefix(path, "/") {
		if index := strings.Index(path, "/sdp-trace/"); index >= 0 {
			return path[index+len("/sdp-trace/"):]
		}
	}
	if strings.HasPrefix(path, "sdp-trace/") {
		return strings.TrimPrefix(path, "sdp-trace/")
	}
	if index := strings.Index(path, "github.com/fall_out_bug/sdp-trace/"); index >= 0 {
		return path[index+len("github.com/fall_out_bug/sdp-trace/"):]
	}
	return strings.TrimPrefix(path, "./")
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(2)
}

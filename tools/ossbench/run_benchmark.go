package main

import (
	"path/filepath"
	"strings"
)

// runBenchmark executes a benchmark definition for n iterations.
func runBenchmark(def benchmarkDef, iterations int) benchmarkResult {
	if iterations <= 0 {
		iterations = 20
	}
	if def.Cmd == "" {
		return missingCommandResult(def, iterations)
	}

	times, attempted, lastErr := runIterations(def, iterations)
	argv := append([]string{filepath.Base(def.Cmd)}, def.Args...)
	cmdDisplay := strings.Join(argv, " ")
	return buildResult(def, cmdDisplay, argv, attempted, times, lastErr)
}

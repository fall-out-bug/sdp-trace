package main

import (
	"fmt"
	"io"
)

func runAndPrint(defs []benchmarkDef, cfg runConfig, stdout, stderr io.Writer) int {
	results, err := runAllBenchmarks(defs, cfg.iterations, cfg.raw)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		return 2
	}
	return printAndExit(stdout, stderr, results, cfg.asJSON)
}

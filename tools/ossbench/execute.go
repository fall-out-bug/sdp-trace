package main

import (
	"fmt"
	"io"
)

func executeBenchmarks(cfg runConfig, stdout, stderr io.Writer) int {
	defs, cleanup, err := resolveBenchmarkDefs(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		return 2
	}
	if cleanup != nil {
		defer cleanup()
	}
	return runAndPrint(defs, cfg, stdout, stderr)
}

package main

import (
	"context"
	"fmt"
	"io"
)

func runWitness(_ context.Context, args []string, stdout, stderr io.Writer) int {
	options, ok := parseWitnessOptions(args, stderr)
	if !ok {
		return exitUsage
	}
	// Witness output is generated from explicit CLI inputs; missing trust
	// material is rejected before record construction.
	record, err := buildWitnessRecord(options)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return writeWitnessRecordOutput(stdout, record)
}

package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func loadProtectedWitnessExpectation(target string, stderr io.Writer) (demo.WitnessExpectation, int, bool) {
	// The expectation loader derives the run id and artifact digests that the
	// supplied witness summary must bind to.
	expected, err := demoWitnessExpectation(target)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return demo.WitnessExpectation{}, exitCannotVerify, false
	}
	return expected, 0, true
}

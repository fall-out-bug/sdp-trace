package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func runStandardGate(target, outPath string, opts *flagSet, stdout, stderr io.Writer) int {
	// Standard gates derive local/CI/audit verdicts from demo run evidence and
	// optional witness data; they do not grant protected-profile trust.
	result, err := demo.WriteGate(target, outPath, opts.stringValue("contract"), opts.stringValue("witness"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return gateExitCode(result)
}

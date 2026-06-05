package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

func writeWitnessRecordOutput(stdout io.Writer, record witness.Record) int {
	// The witness package has already written the record; stdout repeats the
	// generated artifact in a human-reviewable form.
	payload, _ := json.MarshalIndent(record, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	switch record.Status {
	case witness.StatusCannotVerify, witness.StatusNotAssessed:
		// Automation needs cannot_verify/not_assessed separated from ordinary
		// failed witness checks.
		return exitCannotVerify
	case witness.StatusFail:
		return 1
	default:
		return 0
	}
}

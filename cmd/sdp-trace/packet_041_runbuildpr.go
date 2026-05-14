package main

import (
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func runPacketBuildPR(args []string, stdout, stderr io.Writer) int {
	// The CLI layer only parses command intent; packet trust decisions happen
	// after options are converted into a portable evidence input.
	opts, code, ok := parsePacketBuildPROptions(args, stderr)
	if !ok {
		return code
	}
	// build-pr first reconstructs a portable GitHub evidence input; no packet
	// artifact is written until that input and generated bundle validate.
	input, err := buildPRInputFromOptions(opts)
	if err != nil {
		// Input reconstruction failures are emitted as structured cannot_verify
		// output so automation can consume the failed packet attempt.
		result := packet.BuildPRResult{State: packet.StateCannotVerify, Errors: []string{err.Error()}}
		writeJSONPayloadUnchecked(stdout, result)
		return exitCannotVerify
	}
	result, bundle := buildPacketPRResult(input, opts.stringValue("out"))
	if result.State != packet.StatePass {
		// Failed packet gates are surfaced as JSON but not persisted as a
		// complete packet artifact set.
		writeJSONPayloadUnchecked(stdout, result)
		return exitCannotVerify
	}
	return writePacketPRArtifacts(opts.stringValue("out"), bundle, result, stdout, stderr)
}

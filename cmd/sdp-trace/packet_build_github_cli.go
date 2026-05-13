package main

import (
	"fmt"
	"io"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func runPacketBuildGitHub(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parsePacketBuildGitHubOptions(args, stderr)
	if !ok {
		return code
	}
	// GitHub packet input is already materialized JSON; this command only
	// converts it into a portable packet bundle.
	input, err := packet.LoadGitHubInput(opts.stringValue("github-input"))
	if err != nil {
		// Input-read failure is reported as cannot_verify, not a packet fail.
		fmt.Fprintf(stderr, "read github input: %v\n", err)
		return exitCannotVerify
	}
	// Build and validate use a single clock tick for this local conversion.
	bundle := packet.BuildFromGitHubInput(input, time.Now().UTC())
	result := packet.Validate(bundle, time.Now().UTC())
	if result.State != packet.StatePass {
		// Invalid generated bundles are emitted as structured validation evidence.
		writeJSONPayloadUnchecked(stdout, result)
		return exitCannotVerify
	}
	// A passing validation result is the authority to persist the bundle.
	// Only validated bundles are written as durable packet artifacts.
	return writePacketBundle(opts.stringValue("out"), bundle, stdout, stderr)
}

func parsePacketBuildGitHubOptions(args []string, stderr io.Writer) (*flagSet, int, bool) {
	return parsePacketRequiredOptions(args, stderr, "packet build-github", "packet build-github accepts only flags", packetBuildGitHubRequiredFlags)
}

func writePacketBundle(outPath string, bundle packet.Bundle, stdout, stderr io.Writer) int {
	if err := writeJSONFile(outPath, bundle); err != nil {
		// Packet publication requires durable JSON, not stdout-only success.
		fmt.Fprintf(stderr, "write packet bundle: %v\n", err)
		return exitCannotVerify
	}
	fmt.Fprintf(stdout, "wrote %s\n", outPath)
	return 0
}

package main

import (
	"fmt"
	"io"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func runPacketValidate(args []string, stdout, stderr io.Writer) int {
	// Validate mode is intentionally narrow: one bundle in, one verdict out.
	opts := &flagSet{name: "packet validate"}
	opts.setString("bundle", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	// Bundle validation accepts only explicit artifact paths.
	if !requireOnlyFlags(opts, stderr, "packet validate accepts only flags", packetValidateRequiredFlags) {
		return exitUsage
	}
	// Validation reads a committed bundle and writes the verdict JSON unchanged.
	// It does not rebuild packet content from ambient repository state.
	bundle, err := packet.LoadBundle(opts.stringValue("bundle"))
	if err != nil {
		// An unreadable bundle means the packet verdict cannot be replayed.
		fmt.Fprintf(stderr, "read packet bundle: %v\n", err)
		return exitCannotVerify
	}
	// Validation time is local observation metadata for this CLI invocation.
	result := packet.Validate(bundle, time.Now().UTC())
	// Always publish the structured packet verdict before mapping the exit code.
	// Consumers should inspect the JSON state, not infer details from status.
	// The status code is only a shell-friendly summary of that verdict.
	writeJSONPayloadUnchecked(stdout, result)
	return packetValidationExit(result)
}

func runPacketCheckDemo(args []string, stdout, stderr io.Writer) int {
	// Demo-check mode applies the repository's first-packet readiness contract.
	opts := &flagSet{name: "packet check-demo"}
	opts.setString("bundle", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	// Demo checks are flag-only so no untracked positional input is trusted.
	if !requireOnlyFlags(opts, stderr, "packet check-demo accepts only flags", packetCheckDemoRequiredFlags) {
		return exitUsage
	}
	// The demo gate is a stricter first-packet contract over the same bundle.
	// It never consults live GitHub or CI state.
	bundle, err := packet.LoadBundle(opts.stringValue("bundle"))
	if err != nil {
		// Demo checks cannot infer trust from a missing packet bundle.
		fmt.Fprintf(stderr, "read packet bundle: %v\n", err)
		return exitCannotVerify
	}
	// Demo check time is local observation metadata for this CLI invocation.
	result := packet.CheckDemoFirstPacket(bundle, time.Now().UTC())
	// Demo gate output is intentionally the same validation envelope shape.
	// The exit code only distinguishes pass from expected demo-gate failure.
	// The JSON payload remains the detailed evidence for reviewers.
	writeJSONPayloadUnchecked(stdout, result)
	return packetDemoGateExit(result)
}

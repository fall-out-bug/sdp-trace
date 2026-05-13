package main

import (
	"io"
)

func runPreviewCommand(commandName, mode string, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	// Preview identifies the hypothetical command and contract before rendering
	// a non-evidence plan.
	opts, command, code, ok := parsePreviewCommandArgs(commandName, args, stderr)
	if !ok {
		return code
	}
	// Preview loads the contract but deliberately avoids recorder.Run so no run
	// artifacts or trace events are written.
	// Contract parsing is the only validation preview performs.
	contract, code, ok := loadPreviewContract(commandName, opts, stderr)
	if !ok {
		return code
	}
	payload := previewCommandPayload(mode, command, contract)
	writePreviewCommandPayload(stdout, payload)
	return 0
}

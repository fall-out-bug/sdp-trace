package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func parsePreviewCommandArgs(commandName string, args []string, stderr io.Writer) (*flagSet, []string, int, bool) {
	opts := &flagSet{name: commandName}
	// Preview defaults to the built-in contract unless the caller supplies an
	// override, but the output still records which contract was selected.
	opts.setString("contract", "")
	opts.setBool("use-default-contract", true)
	opts.setString("name", "")
	if err := opts.parse(args); err != nil {
		return nil, nil, exitUsage, false
	}
	// Preview command payload is required because the output describes the
	// command that would be observed.
	command := opts.rest()
	if len(command) == 0 {
		// Preview still needs a command descriptor even though it will not run.
		fmt.Fprintf(stderr, "%s requires a command\n", commandName)
		return nil, nil, exitUsage, false
	}
	if missingRequiredContract(opts) {
		// Default contract use is explicit in preview output so dry-run reports
		// remain replayable.
		fmt.Fprintf(stderr, "%s requires --contract unless --use-default-contract is set\n", commandName)
		return nil, nil, exitUsage, false
	}
	// Successful parsing only produces a plan; recorder execution is not
	// reachable from this command path.
	// The command slice is retained only as a descriptor input for the preview
	// payload.
	return opts, command, 0, true
}

func loadPreviewContract(opts *flagSet, stderr io.Writer) (trace.Contract, int, bool) {
	contractPath := opts.stringValue("contract")
	contract := trace.DefaultContract
	if contractPath != "" {
		// A malformed preview contract is cannot_verify because the preview
		// cannot describe valid evidence requirements.
		loaded, err := trace.LoadContract(contractPath)
		if err != nil {
			fmt.Fprintf(stderr, "failed to load contract: %v\n", err)
			return trace.Contract{}, exitCannotVerify, false
		}
		contract = loaded
	}
	return contract, 0, true
}

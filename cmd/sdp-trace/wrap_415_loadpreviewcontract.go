package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func loadPreviewContract(commandName string, opts *flagSet, stderr io.Writer) (trace.Contract, int, bool) {
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

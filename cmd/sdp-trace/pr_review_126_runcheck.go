package main

import (
	"io"
)

func runPRReviewCheck(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parsePRReviewCheckArgs(args, stderr)
	if !ok {
		return code
	}
	// The combined check preserves the packet -> runs -> ledger -> validation
	// sequence so later artifacts can cite earlier ones.
	packet, profile, code, ok := preparePRReviewCheck(opts, args, stderr)
	if !ok {
		return code
	}
	runs, preview, code, ok := executePRReviewCheck(packet, profile, opts, args, stderr)
	if !ok {
		return code
	}
	// Finish handles preview and non-preview publication from the same prepared
	// packet/profile pair.
	return finishPRReviewCheck(opts.stringValue("out"), packet, profile, runs, preview, stdout, stderr)
}

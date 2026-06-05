package main

import (
	"fmt"
	"io"
)

func runPRReviewPacket(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parsePRReviewPacketArgs(args, stderr)
	if !ok {
		return code
	}
	// Packet build errors mean the review input cannot be reconstructed as
	// evidence, so they lower trust instead of becoming a generic CLI failure.
	packet, err := buildPRReviewPacket(opts, args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	writeIndentedPayload(stdout, packet)
	return 0
}

package main

import (
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func writePacketPRArtifacts(outDir string, bundle packet.Bundle, result packet.BuildPRResult, stdout, stderr io.Writer) int {
	markdown, ok := renderPacketPRMarkdown(bundle, &result)
	if !ok {
		writeJSONPayloadUnchecked(stdout, result)
		return exitCannotVerify
	}
	// Durable artifacts are written only after both bundle validation and
	// markdown rendering have succeeded.
	if !writePacketPRFiles(outDir, bundle, result, markdown, stderr) {
		return exitCannotVerify
	}
	writeJSONPayloadUnchecked(stdout, result)
	return 0
}

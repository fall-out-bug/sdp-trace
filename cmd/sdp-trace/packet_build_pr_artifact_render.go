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

func renderPacketPRMarkdown(bundle packet.Bundle, result *packet.BuildPRResult) (string, bool) {
	markdown, err := packet.RenderMarkdown(bundle)
	if err != nil {
		// Rendering failure downgrades the structured build result instead of
		// leaving callers with stderr-only state.
		result.State = packet.StateCannotVerify
		result.Errors = []string{err.Error()}
		return "", false
	}
	return markdown, true
}

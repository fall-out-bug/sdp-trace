package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

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

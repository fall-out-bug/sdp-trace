package main

import (
	"path/filepath"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func buildPacketPRResult(input packet.GitHubPREvidenceInput, outDir string) (packet.BuildPRResult, packet.Bundle) {
	// Build and validation times are generated together so packet rows and the
	// build result describe the same local replay.
	bundle := packet.BuildFromGitHubInput(input, time.Now().UTC())
	validation := packet.Validate(bundle, time.Now().UTC())
	liveGateErrors := packetBuildPRGateErrors(bundle)
	// Result paths are declared before writes so downstream tools can compare
	// the manifest against actual artifact publication.
	result := packet.BuildPRResult{
		State:      packet.StatePass,
		BundlePath: filepath.Join(outDir, "bundle.json"),
		PacketPath: filepath.Join(outDir, "change-evidence-packet.md"),
		ResultPath: filepath.Join(outDir, "build-pr-result.json"),
		Errors:     append(validation.Errors, liveGateErrors...),
	}
	if len(result.Errors) > 0 {
		// Validation and live-gate defects both lower the build to cannot_verify.
		result.State = packet.StateCannotVerify
	}
	return result, bundle
}

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func writePacketPRFiles(outDir string, bundle packet.Bundle, result packet.BuildPRResult, markdown string, stderr io.Writer) bool {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		// Without the output directory, none of the packet artifacts are durable.
		fmt.Fprintf(stderr, "create packet output dir: %v\n", err)
		return false
	}
	// All packet output files share one directory so downstream PR review can
	// cite a single artifact root.
	return writePacketPRArtifactFiles(packetPRArtifactFiles(bundle, result, markdown), stderr)
}

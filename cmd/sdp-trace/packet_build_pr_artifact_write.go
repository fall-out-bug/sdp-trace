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

func writePacketPRArtifactFiles(files []packetPRArtifactFile, stderr io.Writer) bool {
	for _, file := range files {
		// Stop at the first write failure to avoid publishing a partial packet set
		// as if it were complete.
		if !writePacketPRFile(file, stderr) {
			return false
		}
	}
	return true
}

func writePacketPRFile(file packetPRArtifactFile, stderr io.Writer) bool {
	if err := file.write(); err != nil {
		// Labels name the artifact role without exposing full write internals.
		fmt.Fprintf(stderr, "%s: %v\n", file.label, err)
		return false
	}
	return true
}

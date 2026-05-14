package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func packetPRArtifactFiles(bundle packet.Bundle, result packet.BuildPRResult, markdown string) []packetPRArtifactFile {
	// The result carries the paths that this file list materializes.
	return []packetPRArtifactFile{
		{label: "write packet bundle", write: func() error { return writeJSONFile(result.BundlePath, bundle) }},
		{label: "write packet markdown", write: func() error { return writeTextFileAtomic(result.PacketPath, markdown) }},
		{label: "write packet result", write: func() error { return writeJSONFile(result.ResultPath, result) }},
	}
}

package main

import (
	"fmt"
	"io"
)

func writePacketPRFile(file packetPRArtifactFile, stderr io.Writer) bool {
	if err := file.write(); err != nil {
		// Labels name the artifact role without exposing full write internals.
		fmt.Fprintf(stderr, "%s: %v\n", file.label, err)
		return false
	}
	return true
}

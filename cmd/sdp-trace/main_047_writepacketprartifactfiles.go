package main

import (
	"io"
)

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

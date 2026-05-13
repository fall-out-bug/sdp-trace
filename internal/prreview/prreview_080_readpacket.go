package prreview

import (
	"os"

	"path/filepath"
)

func ReadPacket(path string) (Packet, error) {
	// ReadPacket keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	var packet Packet
	if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {

		path = filepath.Join(path, "packet.json")
	}
	return packet, readJSON(path, &packet)
}

package main

import (
	"os"
	"path/filepath"
)

func packetDir(path string) string {
	// Callers may pass either the packet directory or packet.json path.
	// Runner prompt evidence needs the directory that owns copied refs.
	// A missing path falls back to filepath.Dir so load errors stay explicit.
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return path
	}
	return filepath.Dir(path)
}

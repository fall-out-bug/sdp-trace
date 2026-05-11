package main

import (
	"os"
	"path/filepath"
)

func packetDir(path string) string {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return path
	}
	return filepath.Dir(path)
}

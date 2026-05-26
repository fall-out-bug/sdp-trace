package main

import (
	"bytes"
	"os"
	"path/filepath"
)

func containsTrackedFile(root, f, needle string) bool {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(f)))
	return err == nil && bytes.Contains(data, []byte(needle))
}

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

func sha256File(dir, name string) (string, error) {
	// Digest calculation reads the artifact exactly as retained on disk.
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

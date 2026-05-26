package main

import (
	"bufio"
	"os"
	"path/filepath"
)

func actionPinFindings(root, f string) []string {
	file, err := os.Open(filepath.Join(root, filepath.FromSlash(f)))
	if err != nil {
		return []string{unreadableActionFinding(f)}
	}
	defer file.Close()
	return scanActionPins(bufio.NewScanner(file), f)
}

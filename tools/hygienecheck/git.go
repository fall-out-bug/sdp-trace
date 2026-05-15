package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
)

// gitLsFiles returns the paths currently tracked by git under root.
// It shells out to git because the tool must check the live index, not
// a static manifest.
func gitLsFiles(root string) ([]string, error) {
	cmd := exec.Command("git", "ls-files")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-files: %w", err)
	}
	var files []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		files = append(files, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan git ls-files: %w", err)
	}
	return files, nil
}

package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func commandHelp() (string, error) {
	// The command surface is owned by the compiled CLI help text, so the docs
	// check shells out through Go instead of duplicating command metadata here.
	// Tests execute from this package directory, while CI executes from the
	// checkout root, so the command pins its working directory to the repo root.
	cmd := exec.Command("go", "run", "./cmd/sdp-trace", "--help")
	cmd.Dir = repoRoot()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("run sdp-trace help: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return string(output), nil
}

func repoRoot() string {
	// Resolve from this source file instead of the process working directory so
	// `go test ./tools/doccheck` and `go run ./tools/doccheck` check the same
	// repository files.
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func helpCommands(help string) []string {
	// Help output contains prose and Usage rows; only indented command rows are
	// part of the authoritative command surface.
	commands := []string{"sdp-trace --help"}
	for _, line := range strings.Split(help, "\n") {
		if strings.HasPrefix(line, "  sdp-trace ") {
			commands = append(commands, strings.TrimSpace(line))
		}
	}
	return uniqueSorted(commands)
}

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const agentEntrypoint = "docs/agent-entrypoint.md"

func main() {
	os.Exit(exitCode(run(), os.Stderr))
}

func exitCode(err error, stderr io.Writer) int {
	// Keep this tool as a narrow CI gate: command docs either match live help
	// or the process exits non-zero with the drift details.
	if err == nil {
		return 0
	}
	fmt.Fprintln(stderr, err)
	return 1
}

func run() error {
	// Read live command help and the authoritative command surface in one
	// process so CI checks the current checkout, not a stale generated file.
	help, err := commandHelp()
	if err != nil {
		return err
	}
	docPath := filepath.Join(repoRoot(), agentEntrypoint)
	doc, err := os.ReadFile(docPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", agentEntrypoint, err)
	}
	return compareCommandSurface(help, string(doc))
}

package main

// run.go dispatches between the three schemadoc modes: check, generate, and verify-readme.

import (
	"io"
	"os"
	"path/filepath"
)

// run reads the index and dispatches to the requested operation.
func run(generate, verifyReadme bool) error {
	root := repoRoot()
	idx, err := readIndex(filepath.Join(root, indexPath))
	if err != nil {
		return err
	}

	if generate {
		_, err := io.WriteString(os.Stdout, generateTable(idx))
		return err
	}

	if verifyReadme {
		return checkReadme(root, idx)
	}

	return check(root, idx)
}

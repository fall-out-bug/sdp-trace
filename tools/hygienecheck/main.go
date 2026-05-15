package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// main exits non-zero when the repository contains tracked hygiene violations.
func main() {
	os.Exit(exitCode(run(), os.Stderr))
}

// exitCode translates an error into a process exit code. A nil error means
// success; anything else prints to stderr and returns 1.
func exitCode(err error, stderr io.Writer) int {
	if err == nil {
		return 0
	}
	fmt.Fprintln(stderr, err)
	return 1
}

// run executes the full hygiene check starting from the tool's own directory.
func run() error {
	return runAt(repoRoot())
}

// runAt collects all hygiene findings for the repository at root and returns
// an error listing every violation so CI can surface them together. The
// findings are ordered by the check category that produced them. Each check
// runs independently so one failure does not mask the others.
func runAt(root string) error {
	var findings []string

	tracked, err := gitLsFiles(root)
	if err != nil {
		return err
	}

	findings = append(findings, checkForbiddenTrackedPaths(tracked)...)
	findings = append(findings, checkRootArtifactClutter(tracked)...)
	findings = append(findings, checkRootExecutables(root, tracked)...)
	findings = append(findings, checkAbsoluteHomePaths(root, tracked)...)

	if len(findings) > 0 {
		return fmt.Errorf("hygiene findings:\n%s", strings.Join(findings, "\n"))
	}
	return nil
}

// repoRoot resolves the repository root from the source file location so that
// the tool checks the same tree regardless of the working directory. It falls
// back to the current directory when caller resolution fails.
func repoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

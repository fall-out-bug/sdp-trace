package harnessobs

import (
	"os"
	"path/filepath"
	"strings"
)

func relativeSymlinkTarget(path string) (string, error) {
	// Existing output roots are accepted only after resolving symlinks and
	// re-expressing the target relative to the live working directory.
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(resolved)
	if err != nil {
		return "", err
	}
	return filepath.Rel(cwd, abs)
}

func pathEscapesWorkingDirectory(rel string) bool {
	// filepath.Rel reports outside paths with a leading `..`; absolute paths
	// are also treated as escapes for callers that bypass relative conversion.
	return strings.HasPrefix(rel, "..") || filepath.IsAbs(rel)
}

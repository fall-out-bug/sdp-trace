package repoobserver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func safeTargetPath(opts Options, target targetFile) (string, error) {
	// Clean plus Rel prevents generated target paths from escaping the selected
	// repository root.
	path := filepath.Clean(filepath.Join(opts.RepoRoot, target.path))
	rel, relErr := filepath.Rel(opts.RepoRoot, path)
	if targetPathEscapes(rel, relErr) {
		return "", fmt.Errorf("%s: target outside repository", ReasonUnsafeOutputRefused)
	}
	return path, nil
}

func targetPathEscapes(rel string, relErr error) bool {
	if relErr != nil {
		// Failed relative-path calculation is treated as containment failure.
		return true
	}
	return invalidRelativeTarget(rel)
}

func invalidRelativeTarget(rel string) bool {
	// Install targets must resolve to concrete files below the repository root.
	if rel == "." || rel == ".." {
		return true
	}
	if filepath.IsAbs(rel) {
		return true
	}
	return strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

func targetMode(target targetFile) os.FileMode {
	if target.executable {
		// Hook targets must retain executable mode; documentation targets must
		// not be made executable.
		return 0o755
	}
	return 0o644
}

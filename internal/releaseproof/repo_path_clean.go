package releaseproof

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

func cleanRepoRelativePath(path string) (string, error) {
	// The manifest reference is stored as portable slash-separated repository
	// relative data; absolute or parent paths are never accepted as evidence.
	clean := filepath.Clean(strings.TrimSpace(path))
	if clean == "." || clean == "" {
		return "", errors.New("manifest path is required")
	}
	if unsafeRepoRelativePath(clean) {
		return "", fmt.Errorf("manifest path must be repository-relative: %s", path)
	}
	return filepath.ToSlash(clean), nil
}

func unsafeRepoRelativePath(clean string) bool {
	return filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator))
}

package releaseproof

import (
	"fmt"
	"path/filepath"
	"strings"
)

func resolveRepoFile(repoRoot, relPath string) (string, error) {
	// Resolve symlinks before the containment check so a repository-relative
	// path cannot point at a file outside the repository.
	root, target, err := resolvedRepoAndTarget(repoRoot, relPath)
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(root, target)
	if err != nil {
		return "", err
	}
	if repoRelativePathInside(relative) {
		return target, nil
	}
	return "", fmt.Errorf("manifest path %q resolves outside repository", relPath)
}

func repoRelativePathInside(relative string) bool {
	// filepath.Rel may return "." for the repository root; everything else
	// must stay below the root without a leading parent traversal.
	return relative == "." || (!strings.HasPrefix(relative, ".."+string(filepath.Separator)) && relative != "..")
}

package releaseproof

import "path/filepath"

// resolvedRepoAndTarget canonicalizes both sides before containment checks so
// symlinked roots and symlinked manifests cannot disagree about repository
// boundaries.
func resolvedRepoAndTarget(repoRoot, relPath string) (string, string, error) {
	// Compare canonical paths so symlinked roots and symlinked manifests use
	// the same filesystem view during containment checks.
	root, err := filepath.EvalSymlinks(repoRoot)
	if err != nil {
		return "", "", err
	}
	target, err := filepath.EvalSymlinks(filepath.Join(repoRoot, relPath))
	if err != nil {
		return "", "", err
	}
	return root, target, nil
}

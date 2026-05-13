package releaseproof

import (
	"fmt"
	"os/exec"
	"strings"
)

func RepoRoot(dir string) (string, error) {
	// Git decides the repository root for source-bound proof; callers should not
	// infer it from process working-directory assumptions.
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("repo root cannot be determined from %s: %w", dir, err)
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("repo root cannot be determined from %s", dir)
	}
	return root, nil
}

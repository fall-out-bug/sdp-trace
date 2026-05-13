package trace

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// LocalSourceSnapshot returns a local source digest and state disclosure.
func LocalSourceSnapshot(baseDir string) (string, string) {
	// A clean git tree gets a source-bound digest; dirty or non-git state remains
	// explicitly disclosed instead of upgraded.
	cleanBase := filepath.Clean(baseDir)
	tree, treeErr := gitOutput(cleanBase, "rev-parse", "HEAD^{tree}")
	status, statusErr := gitOutput(cleanBase, "status", "--porcelain")
	if treeErr == nil && statusErr == nil {
		state := "git_tree_clean"
		if strings.TrimSpace(status) != "" {
			state = "git_tree_dirty"
		}
		return SHA256Hex("git-tree:" + strings.TrimSpace(tree)), state
	}
	return SHA256Hex("source-not-assessed:" + cleanBase), "not_assessed"
}

func gitOutput(dir string, args ...string) (string, error) {
	// Git command execution is local repository observation, not remote source
	// proof.
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

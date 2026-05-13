package repoobserver

import (
	"errors"
	"os"
	"path/filepath"
)

func updateGitignore(opts Options) ([]DiffSummary, error) {
	// .gitignore handling is limited to the managed block and never interprets
	// unrelated patterns as proof.
	// Missing .gitignore can be created without force because no user content is
	// overwritten.
	path := filepath.Join(opts.RepoRoot, ".gitignore")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, os.WriteFile(path, []byte(gitignoreBlock), 0o644)
	}
	if err != nil {
		return nil, err
	}
	text := string(data)
	start, end := locateGitignoreBlock(text)
	if start >= 0 {
		return updateSdpTraceGitignoreBlock(opts, path, text, data, start, end)
	}
	return appendGitignoreMarker(path, text)
}

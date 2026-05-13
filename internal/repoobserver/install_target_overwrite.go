package repoobserver

import (
	"fmt"
	"os"
	"path/filepath"
)

func overwriteTarget(target targetFile, path string, mode os.FileMode, existing, data []byte) ([]DiffSummary, error) {
	// Backup first, then write; this gives force mode a local recovery path.
	// The summary stores only digests, byte counts, and line counts for safe
	// review.
	// Directory creation still happens after backup so existing content is
	// protected before the replacement write.
	// The replacement write uses the same target mode as first-time generation.
	if err := os.WriteFile(path+".bak", existing, 0o644); err != nil {
		return nil, fmt.Errorf("%s: backup failed for %s", ReasonUnsafeOutputRefused, target.path)
	}
	summary := DiffSummary{
		Path:    target.path,
		Action:  "overwrite_existing_file",
		Before:  contentSummary(existing),
		After:   contentSummary(data),
		Summary: "replace generated file content using safe byte and line counts",
		Backup:  target.path + ".bak",
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return []DiffSummary{summary}, os.WriteFile(path, data, mode)
}

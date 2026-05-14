package main

import (
	"os"
	"path/filepath"
)

func writeTextFileAtomic(path, value string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	// Text artifacts are written through a sibling temp file so readers never
	// observe a partially rendered Markdown/report file.
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	return finishAtomicTextWrite(tmp, tmpName, path, value)
}

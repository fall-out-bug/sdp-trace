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

func finishAtomicTextWrite(tmp *os.File, tmpName, path, value string) error {
	if err := writeAndCloseTempText(tmp, value); err != nil {
		return err
	}
	// Permissions are normalized before rename so the final artifact has the
	// same readable mode as other generated evidence files.
	if err := os.Chmod(tmpName, 0o644); err != nil {
		return err
	}
	// Rename is the publication step for the completed text artifact.
	return os.Rename(tmpName, path)
}

func writeAndCloseTempText(tmp *os.File, value string) error {
	if _, err := tmp.WriteString(value); err != nil {
		// Close on write failure so the temp file handle is not leaked before
		// caller cleanup removes it.
		_ = tmp.Close()
		return err
	}
	return tmp.Close()
}

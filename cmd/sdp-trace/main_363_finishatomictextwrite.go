package main

import (
	"os"
)

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

package main

import (
	"os"
)

func writeAndCloseTempText(tmp *os.File, value string) error {
	if _, err := tmp.WriteString(value); err != nil {
		// Close on write failure so the temp file handle is not leaked before
		// caller cleanup removes it.
		_ = tmp.Close()
		return err
	}
	return tmp.Close()
}

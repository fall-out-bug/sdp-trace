package prreview

import (
	"os"
)

func dirHasEntries(entries []os.DirEntry, err error) bool {
	return err == nil && len(entries) > 0
}

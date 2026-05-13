package prreview

import (
	"errors"
	"os"
)

func readDirFailed(err error) bool {
	return err != nil && !errors.Is(err, os.ErrNotExist)
}

package harnessobs

import (
	"path/filepath"
)

func absolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

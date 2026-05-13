package harnessobs

import (
	"path/filepath"

	"strings"
)

func unsafePotentialParentPath(path string) bool {
	return filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..")
}

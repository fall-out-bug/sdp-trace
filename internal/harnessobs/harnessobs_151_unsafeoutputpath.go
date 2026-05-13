package harnessobs

import (
	"path/filepath"

	"strings"
)

func unsafeOutputPath(path string) bool {
	return filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..")
}

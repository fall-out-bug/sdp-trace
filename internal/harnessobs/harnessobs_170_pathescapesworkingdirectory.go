package harnessobs

import (
	"path/filepath"

	"strings"
)

func pathEscapesWorkingDirectory(rel string) bool {
	return strings.HasPrefix(rel, "..") || filepath.IsAbs(rel)
}

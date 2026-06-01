package harnessobs

import (
	"path/filepath"
	"strings"
)

func unsafeEventRefPath(ref string) bool {
	return strings.Contains(ref, "\\") || strings.Contains(ref, "..") || filepath.IsAbs(ref)
}

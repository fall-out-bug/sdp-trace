package witness

import (
	"os"
	"path/filepath"
	"strings"
)

func unsafeInputPath(path string) bool {
	// Reject unsafe text first, then reject symlinks so Customer PKI inputs stay
	// source-bound to explicit files.
	lower := strings.ToLower(filepath.ToSlash(path))
	if unsafeInputPathText(path, lower) {
		return true
	}
	return inputPathIsSymlink(path)
}

func unsafeInputPathText(path, lower string) bool {
	return emptyOrNULPath(path) ||
		unsafeLowerInputPathText(lower)
}

func emptyOrNULPath(path string) bool {
	return strings.TrimSpace(path) == "" || strings.Contains(path, "\x00")
}

func unsafeLowerInputPathText(lower string) bool {
	return strings.Contains(lower, "://") || strings.Contains(lower, "..") || strings.Contains(lower, "private.key")
}

func inputPathIsSymlink(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode()&os.ModeSymlink != 0
}

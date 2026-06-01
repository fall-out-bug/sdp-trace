package harnessobs

import (
	"errors"
	"path/filepath"
	"strings"
)

func safeOutFile(path string) (string, error) {
	// Output files are local artifact targets. The file itself may be missing,
	// but its parent must resolve inside the working directory.
	if unsafeOutputPath(path) {
		return "", errors.New("out must be a relative local file without traversal")
	}

	parent, err := safeParentDir(filepath.Dir(path))
	if err != nil {
		return "", err
	}
	base := filepath.Base(path)
	if !safeOutputBaseName(base) {
		return "", errors.New("unsafe output filename")
	}

	return filepath.Join(parent, base), nil
}

func unsafeOutputPath(path string) bool {
	return filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..")
}

func safeOutputBaseName(base string) bool {
	// The stem must be a safe file identifier while the extension remains
	// allowed; path separators are rejected after filepath.Base normalization.
	stem := strings.TrimSuffix(base, filepath.Ext(base))

	return safeFileIDPattern.MatchString(stem) && !strings.ContainsAny(base, `/\`)
}

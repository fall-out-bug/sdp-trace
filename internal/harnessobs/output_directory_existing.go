package harnessobs

import "errors"

func safeExistingOutDir(clean string) (string, error) {
	// Existing output directories are resolved through symlinks before being
	// accepted as artifact roots.
	rel, err := relativeSymlinkTarget(clean)
	if err != nil {
		return "", err
	}
	if pathEscapesWorkingDirectory(rel) {
		return "", errors.New("out path escapes working directory")
	}

	return ensureOutDirEmptyOrMissing(rel)
}

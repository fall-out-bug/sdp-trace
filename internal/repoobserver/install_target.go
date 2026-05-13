package repoobserver

import (
	"errors"
	"os"
)

func writeTarget(opts Options, target targetFile) ([]DiffSummary, error) {
	// Every generated target is resolved through containment checks before any
	// read or write.
	// Existing files and new files take separate paths so force-mode overwrite
	// policy cannot affect first-time installs.
	path, err := safeTargetPath(opts, target)
	if err != nil {
		return nil, err
	}
	mode := targetMode(target)
	data := []byte(target.content)
	if existing, err := os.ReadFile(path); err == nil {
		return writeExistingTarget(opts, target, path, mode, existing, data)
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return writeNewTarget(path, data, mode)
}

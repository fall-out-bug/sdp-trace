package repoobserver

import (
	"fmt"
	"os"
)

func writeExistingTarget(opts Options, target targetFile, path string, mode os.FileMode, existing, data []byte) ([]DiffSummary, error) {
	// Differing files are overwritten only with --force after backup/diff summary
	// protections are available.
	if string(existing) == target.content {
		if target.executable {
			return nil, os.Chmod(path, mode)
		}
		return nil, nil
	}
	if !opts.Force {
		return nil, fmt.Errorf("%s: %s exists and differs; use --force after reviewing safe diff", ReasonManualStepRequired, target.path)
	}
	return overwriteTarget(target, path, mode, existing, data)
}

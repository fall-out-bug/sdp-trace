package repoobserver

import (
	"errors"
	"os"
	"path/filepath"
)

const ReasonAlreadyInstalled = "already_installed"

func generatedFileSurface(opts Options, rel, surfaceID string) Surface {
	// Generated file presence is an install signal only; proof still requires an
	// observed run or CI artifact.
	path := filepath.Join(opts.RepoRoot, rel)
	_, err := os.Stat(path)
	if err == nil {
		return surface(surfaceID, StatePass, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonAlreadyInstalled, rel, "")
	}
	if errors.Is(err, os.ErrNotExist) {
		return surface(surfaceID, StateFail, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonManualStepRequired, rel, "write generated observer configuration")
	}
	return surface(surfaceID, StateCannotVerify, StateCannotVerify, ScopeLocalStructural, "filesystem:"+rel, ReasonUnsafeOutputRefused, rel, "fix unreadable generated file path")
}

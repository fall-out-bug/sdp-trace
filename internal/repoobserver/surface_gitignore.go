package repoobserver

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func gitignoreSurface(opts Options) Surface {
	// Only the managed sdp-trace marker block is inspected; unrelated ignore
	// rules are outside this surface.
	rel := ".gitignore"
	data, err := os.ReadFile(filepath.Join(opts.RepoRoot, rel))
	if err == nil {
		return gitignoreContentSurface(rel, string(data))
	}
	if errors.Is(err, os.ErrNotExist) {
		return missingGitignoreSurface(rel)
	}
	return surface(SurfaceGitignore, StateCannotVerify, StateCannotVerify, ScopeLocalStructural, "filesystem:.gitignore", ReasonUnsafeOutputRefused, rel, "fix unreadable .gitignore")
}

func gitignoreContentSurface(rel, data string) Surface {
	if strings.Contains(data, "# sdp-trace begin") && strings.Contains(data, "# sdp-trace end") {
		// Marker presence proves only local ignore configuration, not CI proof.
		return surface(SurfaceGitignore, StatePass, StateNotAssessed, ScopeLocalStructural, "filesystem:.gitignore", ReasonAlreadyInstalled, rel, "")
	}
	return missingGitignoreSurface(rel)
}

func missingGitignoreSurface(rel string) Surface {
	return surface(SurfaceGitignore, StateFail, StateNotAssessed, ScopeLocalStructural, "filesystem:.gitignore", ReasonManualStepRequired, rel, "add sdp-trace ignore block")
}

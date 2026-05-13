package repoobserver

import (
	"errors"
	"os"
	"path/filepath"
)

func hookSurface(opts Options, name, surfaceID string) Surface {
	// Hook existence is local structure; proof remains not_assessed until hook
	// output is observed from a git operation.
	rel := filepath.Join(".githooks", name)
	path := filepath.Join(opts.RepoRoot, rel)
	info, err := os.Stat(path)
	if err == nil {
		return presentHookSurface(info, surfaceID, rel)
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return surface(surfaceID, StateCannotVerify, StateCannotVerify, ScopeLocalStructural, "filesystem:"+rel, ReasonUnsafeOutputRefused, rel, "fix unreadable hook path")
	}
	return surface(surfaceID, StateFail, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonHookScriptAbsent, rel, "install generated hook script")
}

func presentHookSurface(info os.FileInfo, surfaceID, rel string) Surface {
	// A non-executable hook is present but cannot be verified as runnable.
	if info.IsDir() {
		return surface(surfaceID, StateFail, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonHookScriptAbsent, rel, "install generated hook script")
	}
	state := StatePass
	if info.Mode()&0o111 == 0 {
		state = StateCannotVerify
	}
	return surface(surfaceID, state, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonHookScriptPresent, rel, "run a git operation to observe hook output")
}

package repoobserver

import "strings"

const (
	ReasonHooksPathAbsent       = "hooks_path_absent"
	ReasonHooksPathSet          = "hooks_path_set"
	ReasonHooksPathMismatch     = "hooks_path_mismatch"
	ReasonHookScriptAbsent      = "hook_script_absent"
	ReasonHookScriptPresent     = "hook_script_present"
	ReasonHookOutputNotObserved = "hook_output_not_observed"
	ReasonLocalHooksBypassable  = "local_hooks_bypassable"
)

func hooksPathSurface(opts Options) Surface {
	// core.hooksPath is checkout-local git config, so it can prove installation
	// shape but not committed repository proof.
	value := strings.TrimSpace(gitOutput(opts.RepoRoot, "config", "--get", "core.hooksPath"))
	if value == ".githooks" {
		return surface(SurfaceHooksPath, StatePass, StateNotAssessed, ScopeLocalStructural, "git_config:core.hooksPath", ReasonHooksPathSet, "core.hooksPath=.githooks", "")
	}
	if value == "" {
		return surface(SurfaceHooksPath, StateFail, StateNotAssessed, ScopeLocalStructural, "git_config:core.hooksPath", ReasonHooksPathAbsent, "", "set core.hooksPath to .githooks")
	}
	return surface(SurfaceHooksPath, StateFail, StateCannotVerify, ScopeLocalStructural, "git_config:core.hooksPath", ReasonHooksPathMismatch, "core.hooksPath="+safeRef(value), "inspect existing hooks path before replacing it")
}

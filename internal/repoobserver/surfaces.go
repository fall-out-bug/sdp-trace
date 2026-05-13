package repoobserver

import "strings"

const (
	SurfaceRepositoryIdentity = "repository_identity"
	SurfaceHooksPath          = "git_hooks_path"
	SurfacePreCommitHook      = "git_hook_pre_commit"
	SurfacePostCommitHook     = "git_hook_post_commit"
	SurfacePrePushHook        = "git_hook_pre_push"
	SurfaceConfig             = "sdp_trace_config"
	SurfaceGitignore          = "sdp_trace_gitignore"

	SurfaceCIWorkflow                  = "github_actions_workflow"
	SurfaceCIArtifactUpload            = "github_actions_artifact_upload"
	SurfaceCIArtifactBundleObservation = "github_actions_artifact_bundle"
	SurfacePRCheckBinding              = "pr_check_binding"
	SurfaceLocalWrappedCommands        = "local_wrapped_commands"
	SurfaceAgentPrompt                 = "agent_prompt"
)

func buildSurfaces(opts Options) []Surface {
	// Stable ordering keeps JSON and human output diffable across repeated runs.
	// Each surface reports a separate trust boundary instead of contributing to
	// an opaque health score.
	// Local hook, config, CI workflow, and non-applicable profile surfaces are
	// listed together so gaps stay visible.
	return []Surface{
		repositoryIdentitySurface(opts),
		hooksPathSurface(opts),
		hookSurface(opts, "pre-commit", SurfacePreCommitHook),
		hookSurface(opts, "post-commit", SurfacePostCommitHook),
		hookSurface(opts, "pre-push", SurfacePrePushHook),
		generatedFileSurface(opts, ".sdp-trace/config.json", SurfaceConfig),
		gitignoreSurface(opts),
		ciWorkflowSurface(opts),
		ciArtifactUploadSurface(opts),
		ciArtifactBundleSurface(opts),
		prCheckBindingSurface(),
		localWrappedCommandsSurface(),
		agentPromptSurface(),
	}
}

func applyInstallPreviewActions(surfaces []Surface) {
	// Preview rewrites only remediation text; it does not change measured
	// install/proof state.
	for i := range surfaces {
		if surfaces[i].InstallState == StateFail && strings.HasPrefix(surfaces[i].NextAction, "run install") {
			surfaces[i].NextAction = "rerun with --write to install this surface"
		}
	}
}

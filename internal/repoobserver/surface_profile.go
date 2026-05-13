package repoobserver

const (
	ReasonAgentReportedNotProof = "agent_reported_not_proof"
	ReasonOutsideProfileScope   = "outside_profile_scope"
)

func agentPromptSurface() Surface {
	// Prompt text is context only; it is never accepted as repository setup
	// proof.
	return surface(SurfaceAgentPrompt, StateNotAssessed, StateNotAssessed, ScopeAgentReported, "agent_prompt:not_inspected", ReasonAgentReportedNotProof, "", "do not rely on prompt instructions as setup proof")
}

func prCheckBindingSurface() Surface {
	// PR check binding belongs to provider evidence, not the local git-hook
	// profile.
	return surface(SurfacePRCheckBinding, StateNotAssessed, StateNotAssessed, ScopeNotApplicable, "github_pr_checks:not_inspected", ReasonOutsideProfileScope, "", "outside selected profile; no action required")
}

func localWrappedCommandsSurface() Surface {
	// Wrapped command observations belong to run evidence, not setup proof.
	return surface(SurfaceLocalWrappedCommands, StateNotAssessed, StateNotAssessed, ScopeNotApplicable, "sdp_trace_runs:not_inspected", ReasonOutsideProfileScope, "", "outside selected profile; no action required")
}

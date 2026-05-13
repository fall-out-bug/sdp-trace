package repoobserver

import "time"

const SchemaVersion = "block28-repo-observer-status-v1"

type Status struct {
	SchemaVersion     string        `json:"schema_version"`
	Profile           string        `json:"profile"`
	RepositoryID      string        `json:"repository_id"`
	RepositoryRootRef string        `json:"repository_root_ref"`
	GitHead           string        `json:"git_head"`
	GitBranch         string        `json:"git_branch"`
	InstallState      string        `json:"install_state"`
	ProofState        string        `json:"proof_state"`
	Surfaces          []Surface     `json:"surfaces"`
	Gaps              []Gap         `json:"gaps"`
	NextActions       []NextAction  `json:"next_actions"`
	ForceDiffSummary  []DiffSummary `json:"force_diff_summary,omitempty"`
	GeneratedAt       string        `json:"generated_at"`
}

func buildStatus(opts Options, installPreview bool) (Status, error) {
	// Derived repository identity is a sanitized local fallback, not externally
	// signed source proof.
	// Status generation is a snapshot; it does not persist files unless Install
	// enters write mode.
	// Git metadata is captured as strings and may be empty when local git cannot
	// provide it.
	// Aggregate states are derived from surface rows after any preview action
	// text is applied.
	// Gaps and next actions are derived from the same surface slice to keep JSON
	// fields consistent.
	// GeneratedAt uses the normalized clock so tests and CLI output can replay
	// deterministic snapshots.
	repoID := opts.RepositoryID
	if repoID == "" {
		repoID = derivedRepositoryID(opts.RepoRoot)
	}
	surfaces := buildSurfaces(opts)
	if installPreview {
		applyInstallPreviewActions(surfaces)
	}
	return statusFromSurfaces(opts, repoID, surfaces), nil
}

func statusFromSurfaces(opts Options, repoID string, surfaces []Surface) Status {
	// Aggregate install/proof states, gaps, and actions are all derived from the
	// same surface snapshot to avoid prose-only closure over missing evidence.
	gitHead, gitBranch := statusGitRefs(opts.RepoRoot)
	state := surfaceStatusState(surfaces)
	return Status{
		SchemaVersion:     SchemaVersion,
		Profile:           opts.Profile,
		RepositoryID:      repoID,
		RepositoryRootRef: "current_repository",
		// Git fields are observations from this checkout, not immutable source
		// proof.
		GitHead:   gitHead,
		GitBranch: gitBranch,
		// Aggregate states are derived, never hand-authored.
		InstallState: state.install,
		ProofState:   state.proof,
		Surfaces:     surfaces,
		// Gaps and actions preserve the same surface evidence used for verdicts.
		Gaps:        state.gaps,
		NextActions: state.actions,
		GeneratedAt: opts.Now.Format(time.RFC3339),
	}
}

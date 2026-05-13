package repoobserver

const (
	// State values intentionally match gate vocabulary while preserving the
	// distinction between local structure and external proof.
	StatePass         = "pass"
	StateFail         = "fail"
	StateNotAssessed  = "not_assessed"
	StateCannotVerify = "cannot_verify"

	// Scope values prevent locally observed files from being reported as CI or
	// external witness evidence.
	ScopeLocalStructural = "local_structural"
	ScopeCIUploaded      = "ci_uploaded"
	ScopeExternalWitness = "external_witnessed"
	ScopeAgentReported   = "agent_reported"
	ScopeNotApplicable   = "not_applicable"
)

func surface(id, install, proof, scope, source, reason, ref, action string) Surface {
	// Sanitization happens at construction so every renderer receives the same
	// safe observed ref.
	return Surface{
		SurfaceID:      id,
		InstallState:   install,
		ProofState:     proof,
		TrustScope:     scope,
		EvidenceSource: source,
		ObservedRef:    safeRef(ref),
		ReasonCode:     reason,
		NextAction:     action,
	}
}

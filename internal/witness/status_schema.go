package witness

const (
	// StatusPass means all required evidence for the selected witness profile
	// was verified.
	StatusPass = "pass"
	// StatusFail means available evidence contradicts the selected witness
	// profile.
	StatusFail = "fail"
	// StatusCannotVerify means required evidence is missing or replay failed.
	StatusCannotVerify = "cannot_verify"
	// StatusNotAssessed keeps an explicit placeholder when a dimension was not
	// evaluated.
	StatusNotAssessed = "not_assessed"
)

const (
	// TrustScopeCIWitnessed requires a CI provider identity that has been bound
	// to independent evidence such as OIDC.
	TrustScopeCIWitnessed = "ci_witnessed"
	// TrustScopeLocalObserved records local environment or filesystem evidence
	// without promoting it to independent CI trust.
	TrustScopeLocalObserved = "local_observed"
	// TrustScopeExternal records evidence from an external authority profile.
	TrustScopeExternal = "external_witnessed"
)

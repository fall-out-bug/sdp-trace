package checkpoint

const (
	CheckpointSchemaVersion     = "block15-signed-checkpoint-v1"
	VerificationSchemaVersion   = "block15-checkpoint-verification-v1"
	PolicySchemaVersion         = "block15-trusted-checkpoint-policy-v1"
	ProfileEd25519Detached      = "sdp-trace-checkpoint/ed25519-detached-v1"
	HashAlgorithmSHA256         = "sha256"
	SignatureAlgorithmEd25519   = "ed25519"
	AuthorityLocalDevelopment   = "local_development_key"
	AuthorityCIIsolatedJob      = "ci_isolated_job"
	AuthorityExternalWitness    = "external_witness_service"
	KeyIsolationNotAssessed     = "not_assessed"
	StatePass                   = "pass"
	StateFail                   = "fail"
	StateCannotVerify           = "cannot_verify"
	StateNotAssessed            = "not_assessed"
	StateNotIntegrated          = "not_integrated"
	TrustScopeLocalSigned       = "local_signed"
	TrustScopeCISigned          = "ci_signed"
	TrustScopeExternalWitnessed = "external_witnessed"
	TrustScopeUntrustedShape    = "untrusted_shape_only"
)

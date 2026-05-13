package adaptercapture

const (
	SchemaVersion = "block19-adapter-capture-assessment-v1"

	ProfileAdapterCapture = "adapter_capture"
	TrustScopeAdapter     = "adapter_capture_observed"
	TrustScopeLocal       = "local_observed"

	StatePass             = "pass"
	StateFail             = "fail"
	StateCannotVerify     = "cannot_verify"
	StateNotAssessed      = "not_assessed"
	StateMissingTelemetry = "missing_telemetry"
	StateNotIntegrated    = "not_integrated"
	StateUnsupported      = "unsupported"
	StateRetentionLimited = "retention_limited"

	BindingSameChain     = "same_chain"
	BindingAdapterBundle = "adapter_bundle"

	IdentitySelfAsserted = "self_asserted"
	IdentityBound        = "bound"

	RetentionDigestOnly          = "digest_only"
	RetentionSanitizedExcerpt    = "sanitized_excerpt"
	RetentionEncryptedRawRef     = "encrypted_raw_ref"
	RetentionExternalArtifactRef = "external_artifact_ref"
	RetentionNotAssessed         = "not_assessed"
)

var adapterConditionIDs = []string{
	"adapter_event_contract_valid",
	"adapter_identity_visible",
	"run_binding_established",
	"task_drift_visible",
	"tool_call_depth_visible",
	"file_mutation_correlated",
	"model_identity_not_overclaimed",
	"test_provenance_not_overclaimed",
	"provider_refs_portable",
	"redaction_metadata_consistent",
	"capture_depth_not_overclaimed",
}

package managed

const (
	SchemaVersion = "block17-managed-assessment-v1"

	ProfileManagedHarness = "managed_harness"
	TrustScopeManaged     = "managed_harness_observed"

	StatePass             = "pass"
	StateFail             = "fail"
	StateCannotVerify     = "cannot_verify"
	StateNotAssessed      = "not_assessed"
	StateMissingTelemetry = "missing_telemetry"
	StateNotIntegrated    = "not_integrated"
	StateUnsupported      = "unsupported"
	StateSuppressed       = "suppressed"

	IdentityVerified     = "verified"
	IdentitySelfClaimed  = "self_claimed"
	IdentityUnauthorized = "unauthorized"
)

var managedConditionIDs = []string{
	"managed_profile_explicitly_selected",
	"managed_policy_loaded",
	"adapter_registry_loaded",
	"managed_boundary_enrolled_before_run",
	"adapter_identity_authorized",
	"adapter_capabilities_satisfy_contract",
	"adapter_activation_observed",
	"adapter_connection_continuous",
	"required_harness_events_observed",
	"required_tool_events_observed",
	"required_file_mutations_observed",
	"test_provenance_not_agent_reported",
	"suppression_policy_valid",
	"bypass_not_observed",
	"managed_witness_bound",
	"override_does_not_upgrade_managed_profile",
}

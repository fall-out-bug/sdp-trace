package ciartifact

const (
	SchemaVersion = "block26-ci-artifact-observation-v1"

	ProfileCIArtifactObservation = "ci_artifact_observation"
	AuthorityScopeObservation    = "ci_artifact_observation"

	StatePass         = "pass"
	StateFail         = "fail"
	StateCannotVerify = "cannot_verify"
	StateNotAssessed  = "not_assessed"

	ProducerCIUploaded          = "ci_uploaded"
	ProducerCheckedIn           = "checked_in"
	ProducerLocalGenerated      = "local_generated"
	ProducerAgentReported       = "agent_reported"
	ProducerHarnessObserved     = "harness_observed"
	ProducerExternalArtifactRef = "external_artifact_ref"
	ProducerNotAssessed         = "not_assessed"

	AccessPresent      = "present"
	AccessAbsent       = "absent"
	AccessPartial      = "partial"
	AccessExpired      = "expired"
	AccessInaccessible = "inaccessible"
	AccessMalformed    = "malformed"
	AccessUnsafe       = "unsafe"
	AccessNotAssessed  = "not_assessed"
	AccessCannotVerify = "cannot_verify"

	BindingMatched       = "matched"
	BindingMismatch      = "mismatch"
	BindingAbsent        = "absent"
	BindingUnverifiable  = "unverifiable"
	BindingNotAssessed   = "not_assessed"
	IndexValid           = "valid"
	IndexSelfReference   = "self_reference"
	IndexDigestMismatch  = "digest_mismatch"
	IndexMissing         = "missing"
	IndexUnverifiable    = "unverifiable"
	IndexNotAssessed     = "not_assessed"
	SafetyRulesetDefault = "block26-default-output-safety-v1"
)

var familyOrder = []string{
	"run",
	"report",
	"witness",
	"provenance",
	"evidence",
	"trace",
	"artifact_index",
	"redaction_scan",
	"review",
	"change_ci",
}

var validFamilies = map[string]bool{
	"run":            true,
	"report":         true,
	"witness":        true,
	"provenance":     true,
	"evidence":       true,
	"trace":          true,
	"artifact_index": true,
	"redaction_scan": true,
	"review":         true,
	"change_ci":      true,
}

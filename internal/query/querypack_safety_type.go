package query

type QueryPackOutputSafety struct {
	VerifiedAbsentSensitiveClasses []string `json:"verified_absent_sensitive_classes"`
	RedactionPolicyDigest          string   `json:"redaction_policy_digest,omitempty"`
}

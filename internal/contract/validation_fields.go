package contract

// Required string fields are grouped by surface so validation failures are
// stable: identity and locking first, then policy fields after event sets.
var contractHeaderFields = []contractStringField{
	{"schema_version", func(c ExpectedEvidenceContract) string { return c.SchemaVersion }},
	{"contract_id", func(c ExpectedEvidenceContract) string { return c.ContractID }},
	{"version", func(c ExpectedEvidenceContract) string { return c.Version }},
	{"contract_source", func(c ExpectedEvidenceContract) string { return c.ContractSource }},
	{"lock_required_before", func(c ExpectedEvidenceContract) string { return c.LockRequiredBefore }},
}

// Event-set fields must be non-empty so an expected-evidence contract cannot
// silently describe a gate with no required observations.
var contractEventSetFields = []contractListField{
	{"required_observer", func(c ExpectedEvidenceContract) []string { return c.RequiredObservers }},
	{"required_event", func(c ExpectedEvidenceContract) []string { return c.RequiredEvents }},
	{"gate_event", func(c ExpectedEvidenceContract) []string { return c.GateEvents }},
}

// Policy fields define how gate trust and retention should be interpreted.
var contractPolicyFields = []contractStringField{
	{"minimum_gate_trust_scope", func(c ExpectedEvidenceContract) string { return c.MinimumGateTrustScope }},
	{"retention_profile", func(c ExpectedEvidenceContract) string { return c.RetentionProfile }},
	{"redaction_profile", func(c ExpectedEvidenceContract) string { return c.RedactionProfile }},
}

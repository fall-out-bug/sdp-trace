package harnessobs

var validRuleKeys = map[string]bool{
	"missing_required_family": true,
	"missing_optional_family": true,
	"source_unavailable":      true,
	"unsafe_input":            true,
	"digest_mismatch":         true,
	"schema_version_mismatch": true,
	"cross_link_conflict":     true,
}

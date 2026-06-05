package harnessobs

var stateRank = map[string]int{
	StateFail:         4,
	StateCannotVerify: 3,
	StateNotAssessed:  2,
	StatePass:         1,
}

var digestFieldNames = map[string]bool{
	"source_digest":     true,
	"validation_digest": true,
	"commit_digest":     true,
	"envelope_digest":   true,
	"payload_digest":    true,
	"sha256":            true,
}

var validFamilies = map[string]bool{

	"harness":     true,
	"model":       true,
	"interaction": true,
	"phase":       true,
	"review":      true,
	"tool":        true,
	"mutation":    true,
	"test":        true,
	"pr":          true,
	"merge":       true,
	"gap":         true,
}

var validContentStates = map[string]bool{
	ContentRedacted:      true,
	ContentDigestOnly:    true,
	ContentRetainedSafe:  true,
	ContentNotApplicable: true,
}

var validStates = map[string]bool{

	StatePass:         true,
	StateFail:         true,
	StateCannotVerify: true,
	StateNotAssessed:  true,
}

var validRuleKeys = map[string]bool{
	"missing_required_family": true,
	"missing_optional_family": true,
	"source_unavailable":      true,
	"unsafe_input":            true,
	"digest_mismatch":         true,
	"schema_version_mismatch": true,
	"cross_link_conflict":     true,
}

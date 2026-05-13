package packet

var theaterReasonCodes = map[string]bool{

	"agent_claimed_verification": true,
	"unbound_intent":             true,
	"ci_theater":                 true,
	"scope_theater":              true,
	"prompt_contamination":       true,
}

// requiredDecisions keeps approval accountability explicit and separate from
// packet evidence organization.

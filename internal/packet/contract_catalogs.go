package packet

const (
	PacketSchemaVersion = "change-evidence-packet.v0"
	BundleSchemaVersion = "evidence-bundle-manifest.v0"

	// Row states are explicit trust states, not numeric health levels.
	StatePass         = "pass"
	StatePartial      = "partial"
	StateFail         = "fail"
	StateCannotVerify = "cannot_verify"
	StateNotAssessed  = "not_assessed"
	StateNotInScope   = "not_in_scope"

	AuthoringToolGenerated = "tool_generated"
	ProjectionCanonical    = "canonical_markdown_artifact"
)

// RequiredRows is the fixed packet contract; unknown rows cannot silently
// extend the trust surface.
var RequiredRows = []string{
	"PC-CHANGE",
	"PC-INITIATOR",
	"PC-AGENT-ROUTE",
	"PC-MUTATION",
	"PC-VERIFICATION",
	"PC-REVIEW",
	"PC-AUTHORITY",
	"PC-THEATER",
	"PC-ATTESTATION",
	"PC-DECISION",
	"PC-RESIDUAL-GAPS",
}

var requiredDecisions = []string{"merge", "release", "risk_acceptance", "security_review"}

var states = map[string]bool{
	StatePass:         true,
	StatePartial:      true,
	StateFail:         true,
	StateCannotVerify: true,
	StateNotAssessed:  true,
	StateNotInScope:   true,
}

var missingReasonStates = map[string]bool{

	StatePartial:      true,
	StateFail:         true,
	StateCannotVerify: true,
	StateNotAssessed:  true,
	StateNotInScope:   true,
}

var packetStates = map[string]bool{

	"draft":        true,
	"review_ready": true,
	"reviewed":     true,
	"superseded":   true,
}

var authoringMethods = map[string]bool{
	AuthoringToolGenerated:         true,
	"hand_authored_before_tooling": true,
}

var retainedForms = map[string]bool{

	"raw":          true,
	"redacted":     true,
	"digest_only":  true,
	"external_ref": true,
	"not_retained": true,
}

var redactionStatuses = map[string]bool{

	"not_needed":      true,
	"redacted":        true,
	"digest_only":     true,
	"withheld":        true,
	StateCannotVerify: true,
}

var theaterReasonCodes = map[string]bool{

	"agent_claimed_verification": true,
	"unbound_intent":             true,
	"ci_theater":                 true,
	"scope_theater":              true,
	"prompt_contamination":       true,
}

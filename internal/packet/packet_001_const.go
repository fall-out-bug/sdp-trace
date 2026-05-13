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

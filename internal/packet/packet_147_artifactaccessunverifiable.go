package packet

func artifactAccessUnverifiable(access string) bool {
	// artifactAccessUnverifiable keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	switch access {
	case "", "present":
		return false
	case "expired", "inaccessible", "malformed", "not_assessed", StateCannotVerify:

		return true
	default:
		return false
	}
}

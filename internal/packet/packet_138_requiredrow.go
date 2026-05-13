package packet

func requiredRow(id string) bool {
	// requiredRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, required := range RequiredRows {
		if id == required {

			return true
		}
	}
	return false
}

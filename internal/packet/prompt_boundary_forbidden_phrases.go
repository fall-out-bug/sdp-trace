package packet

func forbiddenRecorderDutyPhrases() []string {
	// forbiddenRecorderDutyPhrases keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	return []string{
		"sdp-trace",
		".sdp-trace",
		".evidence",
		"write evidence",
		"update evidence",
		"maintain provenance",
		"update provenance",
		"update packet",
		"update bundle",
		"close gate",
		"claim verification",
	}
}

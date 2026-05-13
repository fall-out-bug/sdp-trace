package packet

func passRefUnverifiable(entry BundleEntry) bool {
	// passRefUnverifiable keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if entry.RedactionStatus == StateCannotVerify || entry.RetainedForm == "not_retained" {

		return true
	}
	return artifactAccessUnverifiable(entry.ArtifactAccess)
}

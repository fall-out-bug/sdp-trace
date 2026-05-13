package packet

func demoRouteEvidenceObservedOpenCodeGSDMiniMax(entry BundleEntry) bool {
	// demoRouteEvidenceObservedOpenCodeGSDMiniMax keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if entry.EvidenceKind != "harness_route_observation" || syntheticEntryDigest(entry) {
		return false
	}

	return hasOpenCodeGSDMiniMax(entry.ObservedComponents)
}

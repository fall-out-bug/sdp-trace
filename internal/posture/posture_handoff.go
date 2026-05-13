package posture

// safeHandoff enforces the export handoff boundary. Keys must be safe labels;
// values must not contain unsafe output. Crossing this threshold prevents injection.
func safeHandoff(values map[string]string) bool {
	// safeHandoff keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	for key, value := range values {

		if !safeLabel(key) || unsafeOutput(value) {
			return false
		}
	}
	return true
}

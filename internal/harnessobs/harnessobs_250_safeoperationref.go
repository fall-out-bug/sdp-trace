package harnessobs

func safeOperationRef(ref string) bool {
	// safeOperationRef keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if operationRefPrefix(ref) {

		return safePrefixedOperationRef(ref)
	}
	return safeRef(ref)
}

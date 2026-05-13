package harnessobs

func normalizeAndResolveSessionEventSource(ctx *sessionCollectionContext) (string, error) {
	// normalizeAndResolveSessionEventSource keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	if normalizeErr := normalizeSessionRawEvents(ctx); normalizeErr != nil {
		return "", normalizeErr
	}
	sourcePath, ok := resolvedSessionEventSource(ctx)
	if !ok {
		return "", errSessionSourceUnavailable
	}
	return sourcePath, nil
}

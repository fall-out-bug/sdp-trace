package harnessobs

func resolveMissingSessionEventSource(ctx *sessionCollectionContext) (string, error) {
	// resolveMissingSessionEventSource keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if ctx.profile.RawEventFormat == "" {

		return "", errSessionSourceUnavailable
	}
	return normalizeAndResolveSessionEventSource(ctx)
}

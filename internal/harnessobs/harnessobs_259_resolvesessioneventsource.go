package harnessobs

func resolveSessionEventSource(ctx *sessionCollectionContext) (string, error) {
	// resolveSessionEventSource keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	sourcePath, err := safeProfileRelativeFile(ctx.profilePath, ctx.profile.EventSourcePath)
	if err == nil {
		return sourcePath, nil
	}
	return resolveMissingSessionEventSource(ctx)
}

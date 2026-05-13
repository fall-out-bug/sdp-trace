package harnessobs

func resolvedSessionEventSource(ctx *sessionCollectionContext) (string, bool) {
	// resolvedSessionEventSource keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	sourcePath, err := safeProfileRelativeFile(ctx.profilePath, ctx.profile.EventSourcePath)
	return sourcePath, err == nil
}

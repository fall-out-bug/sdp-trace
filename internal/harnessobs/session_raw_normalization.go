package harnessobs

import "fmt"

// normalizeAndResolveSessionEventSource runs raw normalization and then
// re-validates the configured event source path as the collection input.
func normalizeAndResolveSessionEventSource(ctx *sessionCollectionContext) (string, error) {
	if normalizeErr := normalizeSessionRawEvents(ctx); normalizeErr != nil {
		return "", normalizeErr
	}
	sourcePath, ok := resolvedSessionEventSource(ctx)
	if !ok {
		return "", errSessionSourceUnavailable
	}
	return sourcePath, nil
}

// normalizeSessionRawEvents converts profile-declared raw events into the
// portable event JSONL source and records the normalized source digest.
func normalizeSessionRawEvents(ctx *sessionCollectionContext) error {
	rawPath, err := safeProfileRelativeFile(ctx.profilePath, ctx.profile.RawEventSourcePath)
	if err != nil {
		return fmt.Errorf("raw_event_source_path invalid: %w", err)
	}

	normalizedPath, err := safeProfileRelativeOutFile(ctx.profilePath, ctx.profile.EventSourcePath)
	if err != nil {
		return err
	}
	if err := normalizeRawEvents(ctx.profile.RawEventFormat, rawPath, normalizedPath, sessionCommandFacts(ctx.session), ctx.now); err != nil {
		return err
	}

	ctx.session.NormalizedDigest = digestFile(normalizedPath)
	return nil
}

package harnessobs

import (
	"fmt"
)

func normalizeSessionRawEvents(ctx *sessionCollectionContext) error {
	// normalizeSessionRawEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
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

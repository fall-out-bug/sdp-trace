package harnessobs

import "fmt"

// loadHarnessProfile resolves the harness profile from the session profile's
// directory so collection cannot jump to arbitrary local files.
func loadHarnessProfile(profilePath string, profile SessionProfile) (string, Profile, error) {
	harnessProfilePath, err := safeProfileRelativeFile(profilePath, profile.HarnessProfilePath)
	if err != nil {
		return "", Profile{}, fmt.Errorf("unsafe harness profile path: %w", err)
	}
	harnessProfile, err := LoadProfile(harnessProfilePath)
	if err != nil {
		return "", Profile{}, err
	}
	return harnessProfilePath, harnessProfile, nil
}

// resolveSessionEventSource returns a safe readable event source, normalizing
// configured raw events first when the profile asks for raw capture conversion.
func resolveSessionEventSource(ctx *sessionCollectionContext) (string, error) {
	if ctx.profile.RawEventSourcePath != "" {
		return normalizeAndResolveSessionEventSource(ctx)
	}

	sourcePath, err := safeProfileRelativeFile(ctx.profilePath, ctx.profile.EventSourcePath)
	if err == nil {
		return sourcePath, nil
	}
	return resolveMissingSessionEventSource(ctx)
}

// resolveMissingSessionEventSource distinguishes a missing direct event source
// from a profile that can still produce one through raw-event normalization.
func resolveMissingSessionEventSource(ctx *sessionCollectionContext) (string, error) {
	if ctx.profile.RawEventFormat == "" {
		return "", errSessionSourceUnavailable
	}
	return normalizeAndResolveSessionEventSource(ctx)
}

// resolvedSessionEventSource keeps source resolution tied to the session
// profile directory even after raw normalization has materialized output.
func resolvedSessionEventSource(ctx *sessionCollectionContext) (string, bool) {
	sourcePath, err := safeProfileRelativeFile(ctx.profilePath, ctx.profile.EventSourcePath)
	return sourcePath, err == nil
}

package harnessobs

func validateRequiredSessionPaths(profile SessionProfile) error {
	// validateRequiredSessionPaths keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireNonBlank(profile.HarnessProfilePath, "session profile requires harness_profile_path"); err != nil {
		return err
	}

	if err := requireNonBlank(profile.EventSourcePath, "session profile requires event_source_path"); err != nil {
		return err
	}
	return nil
}

package main

func registerPRReviewCheckFlags(opts *flagSet) {
	// Check mode intentionally mirrors packet and run flags so the one-shot path
	// records the same provenance as the decomposed commands.
	registerPRReviewPacketFlags(opts)
	// Profile, runner policy, and work-dir describe the review boundary; preview
	// selects a dry publication path without changing parsed evidence inputs.
	opts.setString("profile", "")
	opts.setString("allow-external-runner", "")
	opts.setString("work-dir", ".")
	opts.setString("not-assessed-reason", "")
	// Preview changes publication only; it does not add evidence inputs.
	opts.setBool("preview", false)
}

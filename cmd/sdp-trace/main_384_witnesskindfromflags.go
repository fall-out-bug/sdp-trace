package main

func witnessKindFromFlags(opts *flagSet) (string, string, bool) {
	kind := opts.stringValue("kind")
	if !allowedWitnessKind(kind) {
		// The allowed kind list is closed so CLI output maps to known witness
		// schema semantics.
		return "", "witness requires --kind github-actions, gitlab-ci, buildkite, or customer-pki", false
	}
	return kind, "", true
}

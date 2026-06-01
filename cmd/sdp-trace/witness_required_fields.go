package main

func witnessRequiredFieldsFromFlags(opts *flagSet) (witnessRequiredFields, string, bool) {
	// Target is validated before kind-specific flags because every witness
	// record must bind to one observed run root.
	target, message, ok := witnessTargetFromFlags(opts)
	if !ok {
		return witnessRequiredFields{}, message, false
	}
	return witnessKindOutFromFlags(opts, target)
}

func witnessTargetFromFlags(opts *flagSet) (string, string, bool) {
	targets := opts.rest()
	if len(targets) != 1 {
		// One target keeps witness provenance tied to a single run root.
		return "", "witness requires <runs-root-or-run-dir>", false
	}
	return targets[0], "", true
}

func witnessKindOutFromFlags(opts *flagSet, target string) (witnessRequiredFields, string, bool) {
	kind, message, ok := witnessKindFromFlags(opts)
	if !ok {
		return witnessRequiredFields{}, message, false
	}
	// Output path is required for every kind so the witness record is persisted
	// before stdout renders it.
	out, message, ok := witnessOutFromFlags(opts)
	if !ok {
		return witnessRequiredFields{}, message, false
	}
	// Kind-specific validation prevents customer-PKI witnesses from being
	// created without custody/freshness evidence.
	if message, ok := validateWitnessKindFlags(kind, opts); !ok {
		return witnessRequiredFields{}, message, false
	}
	return witnessRequiredFields{target: target, kind: kind, out: out}, "", true
}

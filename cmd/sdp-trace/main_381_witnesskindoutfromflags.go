package main

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

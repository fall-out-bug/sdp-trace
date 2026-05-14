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

package main

func witnessOptionsFromFlags(opts *flagSet) (witnessOptions, string, bool) {
	// Required fields are normalized before optional witness-specific material
	// is copied into the final options struct.
	fields, message, ok := witnessRequiredFieldsFromFlags(opts)
	if !ok {
		return witnessOptions{}, message, false
	}
	return witnessOptionsFromRequiredFields(fields, opts), "", true
}

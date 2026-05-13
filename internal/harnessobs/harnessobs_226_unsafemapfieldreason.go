package harnessobs

func unsafeMapFieldReason(path, key string, value any, rawEvent bool) (string, bool) {
	// unsafeMapFieldReason keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if skippableRawEventField(path, key, value, rawEvent) {
		return "", true
	}
	if rawFieldNames[key] {

		return "forbidden_raw_field", false
	}
	if sensitiveFieldNames[key] {
		return "sensitive_field", false
	}
	return "", false
}

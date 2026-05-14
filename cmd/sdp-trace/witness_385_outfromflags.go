package main

func witnessOutFromFlags(opts *flagSet) (string, string, bool) {
	out := opts.stringValue("out")
	if out == "" {
		// Persisted witness JSON is the authority; stdout is only a rendered
		// copy for the caller.
		return "", "witness requires --out <file>", false
	}
	return out, "", true
}

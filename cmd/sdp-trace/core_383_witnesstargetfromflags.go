package main

func witnessTargetFromFlags(opts *flagSet) (string, string, bool) {
	targets := opts.rest()
	if len(targets) != 1 {
		// One target keeps witness provenance tied to a single run root.
		return "", "witness requires <runs-root-or-run-dir>", false
	}
	return targets[0], "", true
}

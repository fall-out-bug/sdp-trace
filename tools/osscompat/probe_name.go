package main

func canonicalProbeName(name string) string {
	// Keep legacy one-off probe names usable without adding duplicate entries
	// to the default all-probe run.
	if canonical, ok := legacyProbeNames[name]; ok {
		return canonical
	}
	return name
}

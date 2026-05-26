package main

func collectRegistryUsages(surface commandSurfaceSchema) map[string]bool {
	// The registry owns both canonical Usage rows and example Variations, so
	// drift detection must treat both as copy-pasteable public command lines.
	usages := make(map[string]bool)
	for _, cmd := range surface.Commands {
		addCommandSurfaceUsages(usages, cmd)
	}
	return usages
}

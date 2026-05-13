package main

// defaultAnalysisPaths supplies the repository-root default analysis scope.
func defaultAnalysisPaths(paths []string) []string {
	if len(paths) != 0 {
		// Positional arguments are analyzed paths; flags select which gates or
		// renderers consume the measured report.
		return paths
	}
	// Empty path input means repository root by convention; all threshold and
	// baseline defaults stay disabled until their flags are explicitly set.
	return []string{"."}
}

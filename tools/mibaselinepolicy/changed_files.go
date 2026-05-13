package main

func hasProductionGoChange(changed []string) bool {
	for _, path := range changed {
		// Baseline ratchets are protected only when active production Go code
		// changes in the same diff.
		if productionGoFile(path) {
			return true
		}
	}
	return false
}

func changedSet(changed []string) map[string]bool {
	set := make(map[string]bool, len(changed))
	for _, path := range changed {
		// Exact repo-relative paths are used because baseline policy subjects
		// are fixed files, not glob patterns.
		set[path] = true
	}
	return set
}

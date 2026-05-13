package main

func functionKey(fn functionMetric) string {
	// Function baselines intentionally key by normalized file and display name;
	// line changes alone should not orphan an existing ratchet.
	return fn.file + ":" + fn.name
}

func fileKey(file fileMetric) string {
	// File baselines are path-scoped because the full file is the measured
	// subject.
	return file.file
}

package posture

func explainHeaderLines(result ExportResult) []string {
	// explainHeaderLines keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	return []string{
		"schema_version=" + result.SchemaVersion,
		"export_profile_id=" + result.ExportProfileID,
		"grouping_set_id=" + result.GroupingSetID,
	}
}

package posture

func unsupportedExportHeader(result ExportResult) bool {
	return result.SchemaVersion != SchemaVersion ||
		result.ExportProfileID != ProfileID ||
		result.ExportProfileVersion != ProfileVer
}

func malformedExportHeader(result ExportResult) bool {
	return result.ExportID == "" || result.Producer != "sdp-trace" || result.GeneratedAt == ""
}

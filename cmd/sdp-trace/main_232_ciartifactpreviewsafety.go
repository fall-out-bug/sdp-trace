package main

var ciArtifactPreviewSafety = map[string]string{
	"raw_artifact_content": "not_rendered",
	"reason_payloads":      "safe_reason_codes_only",
	"network_fetch":        "not_performed",
}

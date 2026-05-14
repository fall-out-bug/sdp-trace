package main

var adapterCapturePreviewExpectedEvidence = map[string]string{
	"binding_modes":        "same_chain,adapter_bundle",
	"test_provenance":      "ci_executed,wrapper_executed,harness_observed,agent_reported,cannot_verify",
	"capture_depth_states": "captured,missing_telemetry,not_integrated,unsupported,retention_limited,not_assessed,cannot_verify",
}

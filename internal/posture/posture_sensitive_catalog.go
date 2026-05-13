package posture

func SensitiveClasses() []string {
	return append([]string(nil), sensitiveClasses...)
}

var sensitiveClasses = []string{"raw_command_args", "command_name_or_path", "unsafe_test_identifier", "stdout_stderr_body", "prompt_body", "source_snippet", "tool_payload", "adapter_configuration", "gateway_evidence_ref", "credential_or_token", "authenticated_provider_url", "model_request_response_payload", "raw_review_body", "unsafe_raw_reference_note", "private_filesystem_path", "unsafe_personal_identifier", "unsafe_label", "raw_digest_manifest_path", "free_text_exception_or_refusal_reason"}

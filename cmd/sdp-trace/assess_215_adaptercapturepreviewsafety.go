package main

var adapterCapturePreviewSafety = map[string]string{
	"raw_payloads":    "not_rendered",
	"adapter_secrets": "not_rendered",
	"gateway_refs":    "token_free_refs_only",
	"model_payloads":  "digest_or_block18_reference_only",
}

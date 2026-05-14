package main

var authorityPreviewSafety = map[string]string{
	"raw_prompts":       "not_accepted",
	"raw_model_outputs": "not_accepted",
	"credential_refs":   "rejected_as_malformed",
	"policy_effects":    "not_emitted",
}

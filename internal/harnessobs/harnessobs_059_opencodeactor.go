package harnessobs

func openCodeActor(raw map[string]any) string {
	// openCodeActor keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if model := findStringByKey(raw, "model", "model_id", "modelid"); model != "" {

		return safeToken(model)
	}
	if provider := findStringByKey(raw, "provider"); provider != "" {
		return safeToken(provider)
	}
	return "opencode"
}

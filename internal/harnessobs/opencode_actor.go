package harnessobs

func openCodeActor(raw map[string]any) string {
	if model := findStringByKey(raw, "model", "model_id", "modelid"); model != "" {
		return safeToken(model)
	}
	if provider := findStringByKey(raw, "provider"); provider != "" {
		return safeToken(provider)
	}
	return "opencode"
}

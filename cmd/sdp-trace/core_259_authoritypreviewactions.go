package main

func authorityPreviewActions(inputs map[string]string) []string {
	return previewActionForInputState(inputs, "authority_package", "Supply authority package before authority envelope assessment.", "Fix authority package so it is readable JSON.")
}

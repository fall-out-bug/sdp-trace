package main

func adapterCapturePreviewActions(inputs map[string]string) []string {
	return previewActionForInputState(inputs, "run", "Supply run before adapter capture assessment.", "Fix run so it is a readable JSON run directory.")
}

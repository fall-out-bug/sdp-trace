package harnessobs

func openCodeToolFamily(raw map[string]any, signals []string) bool {
	return hasKey(raw, "tool", "tool_call", "toolcall") ||
		hasSignal(signals, "tool.call", "tool.result", "tool_use") ||
		hasSignalPrefix(signals, "tool.")
}

func openCodeMutationFamily(raw map[string]any, signals []string) bool {
	return hasSignal(signals, "file.write", "file.edit", "file.patch", "file.delete", "mutation") ||
		hasSignalPrefix(signals, "mutation.") ||
		nativeMutationTool(raw)
}

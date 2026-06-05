package harnessobs

func extractCommandModel(command []string) string {
	if shellCommand := shellCommandString(command); shellCommand != "" {
		if model := extractCommandModelArgs(shellFields(shellCommand)); model != "" {
			return model
		}
	}
	return extractCommandModelArgs(command)
}

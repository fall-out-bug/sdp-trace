package main

func missingQuickstartCommands(qsCmds []string) []string {
	qsSet := stringSliceToSet(qsCmds)
	var missing []string
	for _, req := range requiredQuickstartCommands {
		if !quickstartHasCommand(qsSet, req) {
			missing = append(missing, req)
		}
	}
	return missing
}

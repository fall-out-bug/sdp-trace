package main

import "fmt"

const contributorQuickstart = "docs/contributor-quickstart.md"

func compareQuickstartWithRegistry(quickstart string, registry []string) error {
	qsCmds := quickstartCommands(quickstart)
	missing := missingQuickstartCommands(qsCmds)
	stale := staleQuickstartCommands(qsCmds, registry)
	if len(missing) == 0 && len(stale) == 0 {
		return nil
	}
	return fmt.Errorf("quickstart drift: %s", quickstartDrift(missing, stale))
}

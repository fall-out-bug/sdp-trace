package repoobserver

type surfaceState struct {
	install string
	proof   string
	gaps    []Gap
	actions []NextAction
}

func surfaceStatusState(surfaces []Surface) surfaceState {
	// Every aggregate field is recomputed from measured surfaces, not from task
	// checkboxes or generated prose.
	return surfaceState{
		install: aggregateInstallState(surfaces),
		proof:   aggregateProofState(surfaces),
		gaps:    gapsFor(surfaces),
		actions: nextActionsFor(surfaces),
	}
}

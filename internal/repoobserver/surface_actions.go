package repoobserver

import "sort"

func nextActionsFor(surfaces []Surface) []NextAction {
	// Stable sorting makes remediation output deterministic for docs and tests.
	// Empty actions are omitted so already-satisfied surfaces do not create
	// misleading follow-up work.
	actions := make([]NextAction, 0)
	for _, s := range surfaces {
		if s.NextAction == "" {
			continue
		}
		actions = append(actions, nextActionForSurface(s))
	}
	sort.SliceStable(actions, func(i, j int) bool {
		return actions[i].SurfaceID < actions[j].SurfaceID
	})
	return actions
}

func nextActionForSurface(s Surface) NextAction {
	return NextAction{SurfaceID: s.SurfaceID, ActionText: s.NextAction, Blocking: surfaceActionBlocking(s)}
}

func surfaceActionBlocking(s Surface) bool {
	return s.InstallState == StateFail || s.InstallState == StateCannotVerify || s.ProofState == StateCannotVerify
}

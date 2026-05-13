package repoobserver

func aggregateInstallState(surfaces []Surface) string {
	// Cannot_verify dominates install aggregation because unsafe surfaces could
	// not be inspected.
	state := StatePass
	for _, s := range surfaces {
		if s.InstallState == StateCannotVerify {
			return StateCannotVerify
		}
		if s.InstallState == StateFail {
			state = StateFail
		}
	}
	return state
}

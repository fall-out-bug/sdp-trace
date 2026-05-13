package harnessobs

func setFamily(families map[string]bool, family string, observed bool) {
	if observed {
		// Only observed families are recorded; absence remains explicit later.
		families[family] = true
	}
}

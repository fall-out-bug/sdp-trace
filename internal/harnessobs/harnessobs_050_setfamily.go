package harnessobs

func setFamily(families map[string]bool, family string, observed bool) {
	if observed {
		families[family] = true
	}
}

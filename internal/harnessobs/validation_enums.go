package harnessobs

func validFamily(family string) bool {
	return validFamilies[family]
}

func validState(state string) bool {
	return validStates[state]
}

func validContentState(state string) bool {
	return validContentStates[state]
}

func validRuleKey(key string) bool {
	return validRuleKeys[key]
}

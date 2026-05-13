package harnessobs

func validDegradationRule(rule Rule) bool {
	return validState(rule.State) && safeIDPattern.MatchString(rule.ReasonCode)
}

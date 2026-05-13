package query

func familyForEvent(eventType string) string {
	return firstMatchingFamily(eventType, eventFamilyRules, EvidenceFamilyRunChain)
}

func familyForVerifierState(id string) string {
	return firstMatchingFamily(id, verifierFamilyRules, EvidenceFamilyRunChain)
}

package query

import "strings"

type familyRule struct {
	token  string
	family string
}

var eventFamilyRules = []familyRule{
	{"supersed", EvidenceFamilySupersession},
	{"task", EvidenceFamilyTask},
	{"command", EvidenceFamilyCommand},
	{"file", EvidenceFamilyFileMutations},
	{"test", EvidenceFamilyTest},
	{"redaction", EvidenceFamilyRedaction},
}

var verifierFamilyRules = []familyRule{
	{"witness", EvidenceFamilyWitness},
	{"supersed", EvidenceFamilySupersession},
	{"task", EvidenceFamilyTask},
	{"command", EvidenceFamilyCommand},
	{"file", EvidenceFamilyFileMutations},
	{"test", EvidenceFamilyTest},
	{"redaction", EvidenceFamilyRedaction},
}

func firstMatchingFamily(value string, rules []familyRule, fallback string) string {
	// firstMatchingFamily keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	for _, rule := range rules {
		if strings.Contains(value, rule.token) {
			return rule.family
		}
	}
	return fallback
}

package ciartifact

import "sort"

func requiredFamilyObservations(reqs map[string]FamilyRequirement, observed map[string]FamilyInput) ([]FamilyObservation, map[string]bool) {
	// Required family observations preserve the requested producer scope and access
	// state before observed artifacts are matched.

	seen := map[string]bool{}
	out := make([]FamilyObservation, 0, len(reqs)+len(observed))
	for _, family := range familyOrder {
		if req, ok := reqs[family]; ok {
			out = append(out, evaluateFamily(req, observed[family], true))
			seen[family] = true
		}
	}
	return out, seen
}

func extraFamilies(observed map[string]FamilyInput, seen map[string]bool) []string {
	// Extra family detection reports unexpected artifact groups without treating
	// them as proof that required families were covered.
	var extra []string
	for family := range observed {
		if !seen[family] {

			extra = append(extra, family)
		}
	}
	sort.Strings(extra)
	return extra
}

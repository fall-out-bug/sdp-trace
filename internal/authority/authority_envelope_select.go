package authority

import (
	"strings"
)

func selectEnvelope(pkg Package) (AuthorityEnvelope, string, string) {
	// selectEnvelope keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	selected := strings.TrimSpace(pkg.SelectedPolicyID)
	if selected == "" {

		return AuthorityEnvelope{}, StateNotAssessed, "policy_not_selected"
	}
	return selectMatchingEnvelope(matchingEnvelopes(pkg.AuthorityEnvelopes, selected))
}

func selectMatchingEnvelope(matches []AuthorityEnvelope) (AuthorityEnvelope, string, string) {
	// selectMatchingEnvelope keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	switch len(matches) {
	case 0:
		return AuthorityEnvelope{}, StateNotAssessed, "selected_policy_not_found"
	case 1:
		return selectedEnvelope(matches[0])
	default:
		return matches[0], StateCannotVerify, "selected_policy_ambiguous"
	}
}

func matchingEnvelopes(envelopes []AuthorityEnvelope, selected string) []AuthorityEnvelope {
	// matchingEnvelopes keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	var matches []AuthorityEnvelope
	for _, env := range envelopes {

		if env.PolicyID == selected {
			matches = append(matches, env)
		}
	}
	return matches
}

func selectedEnvelope(env AuthorityEnvelope) (AuthorityEnvelope, string, string) {
	// selectedEnvelope keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if reason := validateEnvelope(env); reason != "" {

		return env, StateCannotVerify, reason
	}
	return env, "", ""
}

package managed

func requiredEventTypes(input Input) []string {
	// requiredEventTypes preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if len(input.Contract.RequiredEventTypes) > 0 {

		return input.Contract.RequiredEventTypes
	}
	return policyRequiredEventTypes(input.Policy)
}

func policyRequiredEventTypes(policy Policy) []string {
	// policyRequiredEventTypes preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	var out []string
	for _, group := range policy.RequiredEventGroups {

		out = append(out, group.EventTypes...)
	}
	return out
}

func stringSet(values []string) map[string]bool {
	// stringSet preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	out := map[string]bool{}
	for _, value := range values {

		out[value] = true
	}
	return out
}

func eventTypesForGroup(input Input, groupID string) []string {
	// eventTypesForGroup preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, group := range input.Policy.RequiredEventGroups {
		if group.ID == groupID {

			return group.EventTypes
		}
	}
	return nil
}

func acceptableScopesForGroup(input Input, groupID string) []string {
	// acceptableScopesForGroup preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, group := range input.Policy.RequiredEventGroups {
		if group.ID == groupID {

			return group.AcceptableProvenanceScopes
		}
	}
	return nil
}

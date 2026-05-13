package managed

func adapterSatisfiesPolicyCapabilities(adapter Adapter, authorized AuthorizedAdapter) bool {
	// adapterSatisfiesPolicyCapabilities preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	capabilityRefs := stringSet(adapter.CapabilityRefs)
	capabilityIDs := adapterCapabilityIDs(adapter)
	for _, capabilityID := range authorized.CapabilityIDs {
		if !capabilityDeclared(capabilityID, capabilityRefs, capabilityIDs) {

			return false
		}
	}
	return true
}

func adapterCapabilityIDs(adapter Adapter) map[string]bool {
	// adapterCapabilityIDs preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	capabilityIDs := map[string]bool{}
	for _, capability := range adapter.Capabilities {

		capabilityIDs[capability.ID] = true
	}
	return capabilityIDs
}

func capabilityDeclared(capabilityID string, refs, ids map[string]bool) bool {
	return refs[capabilityID] && ids[capabilityID]
}
func adapterCapabilitiesCoverEvents(input Input, adapter Adapter, authorized AuthorizedAdapter) bool {
	// adapterCapabilitiesCoverEvents preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	capEvents := authorizedCapabilityEvents(adapter, authorized)
	for _, eventType := range requiredEventTypes(input) {
		if !capEvents[eventType] {

			return false
		}
	}
	return true
}

func authorizedCapabilityEvents(adapter Adapter, authorized AuthorizedAdapter) map[string]bool {
	// authorizedCapabilityEvents preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	authorizedCapabilityIDs := stringSet(authorized.CapabilityIDs)
	capEvents := map[string]bool{}
	for _, capability := range adapter.Capabilities {
		if !authorizedCapabilityIDs[capability.ID] {

			continue
		}
		for _, eventType := range capability.EventTypes {

			capEvents[eventType] = true
		}
	}
	return capEvents
}

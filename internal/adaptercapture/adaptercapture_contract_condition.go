package adaptercapture

func contractCondition(run RunEvidence) Condition {
	// Contract evidence is checked before identity and run binding so an adapter
	// cannot claim capture semantics without the expected event contract.
	if len(run.AdapterEvents) == 0 {
		return cannotVerify("adapter_event_contract_valid", "adapter_events_missing", "adapter event evidence is missing", "Supply same-chain adapter events or an adapter bundle.")
	}

	return contractConditionFromEvents(run.AdapterEvents)
}

func contractConditionFromEvents(events []AdapterEvent) Condition {
	// Event-level contract checks keep missing event families separate from malformed
	// or unsupported adapter evidence.
	seen := map[string]bool{}
	for _, event := range events {
		if adapterEventIsMalformed(event) {

			return fail("adapter_event_contract_valid", "adapter_event_malformed", "adapter event is missing required contract fields", "Emit schema-valid adapter events with producer, adapter, type, and digest fields.")
		}
		if hasDuplicateCorrelationKey(seen, event) {

			return cannotVerify("adapter_event_contract_valid", "conflicting_adapter_events", "multiple adapter events share a correlation key", "Deduplicate adapter events or make conflicts explicit.")
		}
	}
	return pass("adapter_event_contract_valid", "adapter_event_contract_valid", "adapter events match required contract fields")
}

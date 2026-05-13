package trace

func validateEventChainIfRequested(events []Event, requireChain bool) error {
	if !requireChain {
		// Shape-only validation callers intentionally skip contiguous chain proof.
		return nil
	}
	// Chain validation binds event order and hashes when requested.
	return ValidateEventChain(events)
}

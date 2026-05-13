package main

type policyInput struct {
	changed        []string
	baselines      []string
	baselineExists func(string) (bool, error)
}

func checkPolicy(input policyInput) error {
	// Baseline ratchets are only required when production Go changes; docs or
	// examples alone should not force a metric-baseline update.
	if !hasProductionGoChange(input.changed) {
		return nil
	}
	// Convert changed paths once so each configured baseline can query scope
	// without repeatedly scanning the raw input slice.
	return checkChangedBaselines(input, changedSet(input.changed))
}

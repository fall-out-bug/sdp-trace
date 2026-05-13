package authority

func Evaluate(pkg Package) Result {
	// Evaluate keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	env, envState, envReason := selectEnvelope(pkg)
	actions := sortedObservedActions(pkg.ObservedActions)
	bindings := evaluateBindings(pkg.EvidenceBindings, actions)
	evaluations := evaluateActions(pkg, env, envState, envReason, actions, bindings)
	return authorityResult(pkg, actions, bindings, evaluations, envState)
}

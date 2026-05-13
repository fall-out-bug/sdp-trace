package managed

func policyCondition(policy Policy) Condition {
	// Policy evidence is required before adapter authorization can be considered
	// managed rather than local observation.
	if policy.PolicyID == "" {

		return cannotVerify("managed_policy_loaded", "missing_managed_policy", "managed policy is required", "Supply a managed policy anchored before the run.")
	}
	if !preRunProvenance(policy.PolicyProvenance.Source) || policy.PolicyProvenance.Digest == "" {

		return fail("managed_policy_loaded", "post_hoc_policy", "managed policy provenance is not anchored before the run", "Use a VCS, CI, human-signed, or customer policy equivalent policy.")
	}
	return pass("managed_policy_loaded", "managed_policy_loaded", "managed policy is readable and anchored before the run")
}

func registryCondition(registry Registry) Condition {
	// Registry evidence names whether the adapter exists in the declared registry
	// instead of treating any selected adapter as trusted.
	if registry.RegistryID == "" {

		return cannotVerify("adapter_registry_loaded", "missing_adapter_registry", "adapter registry is required", "Supply an adapter registry anchored before the run.")
	}
	if !preRunProvenance(registry.Provenance.Source) || registry.Provenance.Digest == "" {

		return fail("adapter_registry_loaded", "post_hoc_registry", "adapter registry provenance is not anchored before the run", "Use a VCS, CI, human-signed, or customer policy equivalent registry.")
	}
	return pass("adapter_registry_loaded", "adapter_registry_loaded", "adapter registry is readable and anchored before the run")
}

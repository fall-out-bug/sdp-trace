package managed

func boundaryCondition(input Input) Condition {
	// Boundary evidence keeps selected adapter identity and task boundary checks
	// explicit before capability or event coverage is evaluated.
	boundary := input.Run.ManagedBoundaryEnrolled
	if boundary == nil {

		return fail("managed_boundary_enrolled_before_run", "run_not_managed", "selected run has no managed boundary enrollment event", "Run through the managed wrapper or lower the claim.")
	}
	if boundaryNotBeforeLaunch(*boundary, input.Run.ChildLaunch) {

		return fail("managed_boundary_enrolled_before_run", "late_enrollment", "managed boundary enrollment is not before child launch", "Enroll the managed boundary before child harness execution.")
	}
	if boundaryBindingMismatch(*boundary, input) {

		return fail("managed_boundary_enrolled_before_run", "managed_boundary_not_in_chain", "managed boundary event does not bind the selected policy, registry, or run nonce", "Regenerate the run under the selected managed policy and registry.")
	}
	return pass("managed_boundary_enrolled_before_run", "managed_boundary_enrolled", "managed boundary enrollment is in chain before child launch")
}

func boundaryNotBeforeLaunch(boundary ManagedBoundaryEnrolled, launch LaunchEvent) bool {
	return boundary.Sequence >= launch.Sequence || launch.Sequence == 0
}

func boundaryBindingMismatch(boundary ManagedBoundaryEnrolled, input Input) bool {
	return boundary.ManagedPolicyDigest != input.Policy.PolicyProvenance.Digest || boundary.AdapterRegistryDigest != input.Registry.Provenance.Digest || boundary.RunNonce != input.Run.RunNonce
}

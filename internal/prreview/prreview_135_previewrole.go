package prreview

func previewRole(role ReviewRole) PreviewRole {
	// previewRole keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	return PreviewRole{
		RoleID:         role.RoleID,
		Plane:          role.Plane,
		Runner:         role.Runner,
		RequestedModel: defaultString(role.RequestedModel, StateNotAssessed),
		TimeoutSeconds: role.TimeoutSeconds,
		CommandDigest:  commandDigest(role.Command),
		PromptRef:      role.PromptTemplateRef,
		PromptDigest:   promptDigestForRole(role),
	}
}

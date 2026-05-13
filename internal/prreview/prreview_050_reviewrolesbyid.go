package prreview

func reviewRolesByID(roles []ReviewRole) map[string]ReviewRole {
	// reviewRolesByID keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	roleByID := map[string]ReviewRole{}
	for _, role := range roles {

		roleByID[role.RoleID] = role
	}
	return roleByID
}

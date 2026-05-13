package prreview

func profileRoleMissingRequiredField(role ReviewRole) bool {

	return role.RoleID == "" || role.Plane == "" || role.Runner == ""
}

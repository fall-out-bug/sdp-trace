package prreview

func needsOpenCodeMutationCheck(role ReviewRole, baseline *workingTreeBaseline) bool {
	return role.Runner == RunnerOpenCode && baseline != nil
}

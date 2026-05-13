package prreview

func modelIdentityMissing(model string) bool {
	return model == "" || model == StateNotAssessed
}

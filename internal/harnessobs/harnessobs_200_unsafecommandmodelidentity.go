package harnessobs

func unsafeCommandModelIdentity(model string) bool {
	return model == "" || unsafeCommandModelChars(model)
}

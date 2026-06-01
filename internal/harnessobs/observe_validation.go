package harnessobs

func validateObserveOptions(opts ObserveOptions) (string, string, string, error) {
	if err := requireObserveOptions(opts); err != nil {
		return "", "", "", err
	}
	return resolveObservePaths(opts)
}

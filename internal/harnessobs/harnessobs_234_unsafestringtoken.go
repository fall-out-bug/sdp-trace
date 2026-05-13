package harnessobs

func unsafeStringToken(path, value string, rawEvent bool) bool {
	return providerTokenPrefix.MatchString(value) || unsafeEncodedToken(path, value, rawEvent)
}

package harnessobs

func safeEncodedTokenExemption(path, value string, rawEvent bool) bool {
	return digestField(path) || sha256Pattern.MatchString(value) || rawEventPathLikeField(path, rawEvent)
}

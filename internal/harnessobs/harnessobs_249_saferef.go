package harnessobs

func safeRef(ref string) bool {
	return safeIDPattern.MatchString(ref) || sha256Pattern.MatchString(ref)
}

package harnessobs

func digestLine(line []byte) string {
	if canonical, ok := canonicalSourceDigestLine(line); ok {
		return sha256Hex(canonical)
	}
	return sha256Hex(line)
}

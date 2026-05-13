package harnessobs

var validContentStates = map[string]bool{
	ContentRedacted:      true,
	ContentDigestOnly:    true,
	ContentRetainedSafe:  true,
	ContentNotApplicable: true,
}

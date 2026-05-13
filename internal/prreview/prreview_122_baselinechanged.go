package prreview

func baselineChanged(after, before *workingTreeBaseline) bool {
	return after.Digest != before.Digest || after.Count != before.Count
}

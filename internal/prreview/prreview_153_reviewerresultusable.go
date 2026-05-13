package prreview

func reviewerResultUsable(result ReviewerResult) bool {
	// Usable review coverage requires both a positive reviewer status and retained
	// raw output evidence. A hand-authored status without a digest-bound output
	// reference stays unverifiable.
	return reviewerStatusUsable(result.Status) && result.RawOutputRef != nil
}

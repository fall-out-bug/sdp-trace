package prreview

func reviewerStatusUsable(status string) bool {
	return status == StatusFindingsReported || status == StatusNoFindings
}

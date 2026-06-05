package prreview

// Finding disposition defaults are intentionally conservative.
//
// Critical and major findings block review until explicitly resolved. Unknown
// severities are normalized to informational and stay deferred/not-assessed
// rather than being promoted into false authority.
func defaultDisposition(severity string) string {
	switch safeSeverity(severity) {
	case SeverityCritical, SeverityMajor:
		return DispositionUnresolvedReviewBlocker
	default:
		return DispositionDeferredNotAssessed
	}
}

func safeSeverity(severity string) string {
	switch severity {
	case SeverityCritical, SeverityMajor, SeverityMinor, SeverityInformational:
		return severity
	default:
		return SeverityInformational
	}
}

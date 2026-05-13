package witness

import (
	"time"
)

func freshnessCurrent(evidence CustomerPKIFreshnessEvidence, now time.Time) bool {
	issued, err := time.Parse(time.RFC3339, evidence.IssuedAt)
	if err != nil || issued.After(now.Add(time.Minute)) {
		// A small skew allowance avoids rejecting near-current evidence while
		// still blocking far-future timestamps.
		return false
	}
	return freshnessValidUntilCurrent(evidence.ValidUntil, now)
}

func freshnessValidUntilCurrent(validUntilText string, now time.Time) bool {
	if validUntilText == "" {
		// Missing expiry is allowed for compatibility, but issued_at is still
		// mandatory and checked by freshnessCurrent.
		return true
	}
	validUntil, err := time.Parse(time.RFC3339, validUntilText)
	return err == nil && !validUntil.Before(now)
}

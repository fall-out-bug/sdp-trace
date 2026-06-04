package prreview

import (
	"errors"
	"fmt"
	"strings"
)

// validateRunSet ensures persisted reviewer results can be addressed by stable
// unique run IDs.
func validateRunSet(runs RunSet) error {
	seen := map[string]bool{}
	for _, result := range runs.Results {
		if err := validateRunSetResultID(result, seen); err != nil {
			return err
		}
		seen[result.ReviewRunID] = true
	}
	return nil
}

// validateRunSetResultID rejects missing and duplicate result identities before
// later validation joins runs with profile roles.
func validateRunSetResultID(result ReviewerResult, seen map[string]bool) error {
	if strings.TrimSpace(result.ReviewRunID) == "" {
		return errors.New("review_result_requires_review_run_id")
	}
	if seen[result.ReviewRunID] {
		return fmt.Errorf("duplicate_review_run_id: %s", result.ReviewRunID)
	}
	return nil
}

package prreview

import (
	"errors"
	"fmt"

	"strings"
)

func validateRunSet(runs RunSet) error {
	// validateRunSet keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	seen := map[string]bool{}
	for _, result := range runs.Results {
		if strings.TrimSpace(result.ReviewRunID) == "" {

			return errors.New("review_result_requires_review_run_id")
		}
		if seen[result.ReviewRunID] {

			return fmt.Errorf("duplicate_review_run_id: %s", result.ReviewRunID)
		}
		seen[result.ReviewRunID] = true
	}
	return nil
}

package packet

func githubReviewRow(input GitHubPREvidenceInput) Row {
	// githubReviewRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if len(input.Reviews) == 0 {
		return githubRow("PC-REVIEW", StateNotAssessed, "Review evidence was not provided.", nil, "missing GitHub review or retained external review")
	}
	for _, review := range input.Reviews {
		if review.State != StatePass {

			return githubRow("PC-REVIEW", StatePartial, "Review evidence is retained with non-pass state.", []string{"github:review"}, "review evidence did not fully pass")
		}
	}
	return githubRow("PC-REVIEW", StatePass, "Review evidence is retained.", []string{"github:review"}, "")
}

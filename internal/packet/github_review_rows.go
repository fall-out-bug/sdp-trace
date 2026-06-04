package packet

// Review rows keep missing review evidence distinct from retained but non-pass
// review evidence. Both are non-pass, but they have different closure paths.
func githubReviewRow(input GitHubPREvidenceInput) Row {
	if len(input.Reviews) == 0 {
		return githubReviewNotAssessedRow()
	}
	if !allGitHubReviewsPass(input.Reviews) {
		return githubReviewPartialRow()
	}
	return githubReviewPassRow()
}

func allGitHubReviewsPass(reviews []GitHubReview) bool {
	for _, review := range reviews {
		if review.State != StatePass {
			return false
		}
	}
	return true
}

func githubReviewNotAssessedRow() Row {
	return githubRow("PC-REVIEW", StateNotAssessed, "Review evidence was not provided.", nil, "missing GitHub review or retained external review")
}

func githubReviewPartialRow() Row {
	return githubRow("PC-REVIEW", StatePartial, "Review evidence is retained with non-pass state.", []string{"github:review"}, "review evidence did not fully pass")
}

func githubReviewPassRow() Row {
	return githubRow("PC-REVIEW", StatePass, "Review evidence is retained.", []string{"github:review"}, "")
}

package packet

// GitHub PR evidence input can carry resolver fields from hand-authored JSON,
// optional enrichment files, and future adapter output. Packet construction
// keeps relative resolver strings valid because examples and local fixtures use
// them, but URL-shaped values need network-safety checks before they become
// retained evidence references.
//
// The collection checks stay separate by input section so error messages point
// reviewers to the evidence field that needs correction. The URL policy itself
// lives in github_resolver_url_validation.go.

// ValidateGitHubPREvidenceInputResolvers rejects URL-shaped resolver values
// that would point generated packet evidence at unsafe network locations.
func ValidateGitHubPREvidenceInputResolvers(input GitHubPREvidenceInput) error {
	if err := validateCheckURLs(input.Checks); err != nil {
		return err
	}
	if err := validateArtifactResolvers(input.Artifacts); err != nil {
		return err
	}
	if err := validateReviewResolvers(input.Reviews); err != nil {
		return err
	}
	return validateIntegrationResolvers(input.IntegrationActions)
}

func validateCheckURLs(checks []GitHubCheck) error {
	for _, check := range checks {
		if err := validateResolverURL("check url", check.URL); err != nil {
			return err
		}
	}
	return nil
}

func validateArtifactResolvers(artifacts []GitHubArtifact) error {
	for _, artifact := range artifacts {
		if err := validateResolverURL("artifact resolver", artifact.Resolver); err != nil {
			return err
		}
	}
	return nil
}

func validateReviewResolvers(reviews []GitHubReview) error {
	for _, review := range reviews {
		if err := validateResolverURL("review resolver", review.Resolver); err != nil {
			return err
		}
	}
	return nil
}

func validateIntegrationResolvers(actions []IntegrationAction) error {
	for _, action := range actions {
		if err := validateResolverURL("integration resolver", action.Resolver); err != nil {
			return err
		}
	}
	return nil
}

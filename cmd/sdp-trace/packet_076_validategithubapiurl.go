package main

func validateGitHubAPIURL(apiURL, serverURL string) error {
	parsed, err := parseGitHubAPIURL(apiURL)
	if err != nil {
		return err
	}
	// Parse and target validation are split so error text can distinguish syntax
	// from trust-target failures.
	return validateParsedGitHubAPIURL(parsed, apiURL, serverURL)
}

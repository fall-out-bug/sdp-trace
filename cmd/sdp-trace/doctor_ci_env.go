package main

import (
	"strings"
)

func requiredCIWitnessEnvFields() []string {
	// Require both OIDC request fields and workflow identity fields for CI
	// witness construction.
	return []string{
		"ACTIONS_ID_TOKEN_REQUEST_TOKEN",
		"ACTIONS_ID_TOKEN_REQUEST_URL",
		"GITHUB_ACTIONS",
		"GITHUB_ACTOR",
		"GITHUB_JOB",
		"GITHUB_REF",
		"GITHUB_REPOSITORY",
		"GITHUB_RUN_ATTEMPT",
		"GITHUB_RUN_ID",
		"GITHUB_SERVER_URL",
		"GITHUB_SHA",
		"GITHUB_WORKFLOW",
	}
}

func missingEnvFields(env map[string]string, required []string) []string {
	missing := make([]string, 0)
	for _, key := range required {
		if strings.TrimSpace(env[key]) == "" {
			// Trim whitespace so empty exported variables are treated as absent.
			missing = append(missing, key)
		}
	}
	return missing
}

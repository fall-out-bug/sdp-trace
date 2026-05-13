package witness

import "sort"

func missingGitHubIdentity(env map[string]string) []string {
	// GitHub identity requires the provider marker plus source and run fields.
	missing := missingEnvKeys(env, githubIdentityEnvKeys)
	if env["GITHUB_ACTIONS"] != "" && env["GITHUB_ACTIONS"] != "true" {
		missing = append(missing, "GITHUB_ACTIONS=true")
	}
	sort.Strings(missing)
	return missing
}

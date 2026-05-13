package ciartifact

import "sort"

func safeClasses(input []string) []string {
	// safeClasses keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	allowed := map[string]bool{
		"token_like": true, "jwt_token": true, "private_key": true,
		"cloud_credential": true, "provider_token": true,
		"private_artifact_url": true, "private_filesystem_path": true,
		"prompt_or_model_payload": true, "raw_job_log": true,
		"high_entropy_secret": true,
	}
	var out []string
	for _, value := range input {
		if allowed[value] {

			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}

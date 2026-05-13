package main

func missingCIWitnessFields(env map[string]string) []string {
	missing := missingEnvFields(env, requiredCIWitnessEnvFields())
	if env["GITHUB_ACTIONS"] != "" && env["GITHUB_ACTIONS"] != "true" {
		// GitHub exposes the flag as literal true for Actions-backed identity.
		missing = append(missing, "GITHUB_ACTIONS=true")
	}
	return missing
}

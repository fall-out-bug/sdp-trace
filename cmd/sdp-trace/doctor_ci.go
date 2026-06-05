package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func ciWitnessPrerequisiteCheck(env map[string]string) doctorCheck {
	missing := missingCIWitnessFields(env)
	if len(missing) > 0 {
		// Local environments usually cannot produce CI witness evidence.
		return doctorCheck{
			ID:      "ci_witness_prerequisites",
			State:   string(trace.VerdictCannotVerify),
			Reason:  "GitHub Actions identity or OIDC prerequisite is unavailable in this environment",
			Missing: missing,
		}
	}
	// Passing prerequisites only means the environment exposes the fields needed
	// for witness construction; it is not a witness verdict.
	return doctorCheck{
		ID:     "ci_witness_prerequisites",
		State:  "pass",
		Reason: "GitHub Actions identity and OIDC prerequisites are present",
	}
}

func missingCIWitnessFields(env map[string]string) []string {
	missing := missingEnvFields(env, requiredCIWitnessEnvFields())
	if env["GITHUB_ACTIONS"] != "" && env["GITHUB_ACTIONS"] != "true" {
		// GitHub exposes the flag as literal true for Actions-backed identity.
		missing = append(missing, "GITHUB_ACTIONS=true")
	}
	return missing
}

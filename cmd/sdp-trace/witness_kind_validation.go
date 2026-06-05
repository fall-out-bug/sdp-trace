package main

import (
	"fmt"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

func witnessKindFromFlags(opts *flagSet) (string, string, bool) {
	kind := opts.stringValue("kind")
	if !allowedWitnessKind(kind) {
		// The allowed kind list is closed so CLI output maps to known witness
		// schema semantics.
		return "", "witness requires --kind github-actions, gitlab-ci, buildkite, or customer-pki", false
	}
	return kind, "", true
}

func witnessOutFromFlags(opts *flagSet) (string, string, bool) {
	out := opts.stringValue("out")
	if out == "" {
		// Persisted witness JSON is the authority; stdout is only a rendered
		// copy for the caller.
		return "", "witness requires --out <file>", false
	}
	return out, "", true
}

func validateWitnessKindFlags(kind string, opts *flagSet) (string, bool) {
	// Missing kind-specific material is a usage failure, not a generated
	// not_assessed witness record.
	missing := missingWitnessKindFlags(kind, opts)
	if len(missing) > 0 {
		return fmt.Sprintf("customer-pki witness requires %s", strings.Join(missing, ", ")), false
	}
	return "", true
}

func missingWitnessKindFlags(kind string, opts *flagSet) []string {
	if kind != witness.KindCustomerPKI {
		// Non-customer-PKI witnesses do not require customer key custody inputs.
		return nil
	}
	return missingCustomerPKIFlags(opts)
}

func allowedWitnessKind(kind string) bool {
	switch kind {
	case witness.KindGitHubActions, witness.KindGitLabCI, witness.KindBuildkite, witness.KindCustomerPKI:
		// Each allowed kind has an explicit builder and schema contract.
		return true
	default:
		return false
	}
}

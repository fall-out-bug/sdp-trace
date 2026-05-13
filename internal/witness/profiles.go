package witness

import (
	"fmt"
)

func BuildProfile(kind, runsRoot, reportDir string, opts ProfileOptions) (Record, error) {
	if ciEnvelopeProfile(kind) {
		// GitLab and Buildkite profiles are envelope-backed because this package
		// does not obtain provider-native attestations for those systems.
		return BuildCIEnvelopeProfile(kind, runsRoot, reportDir, opts.EnvelopePath)
	}
	switch kind {
	case KindCustomerPKI:
		// Customer PKI is external evidence: the profile can pass only after
		// policy, signer, freshness, signature, and run binding all agree.
		return BuildCustomerPKI(runsRoot, reportDir, opts)
	default:
		return Record{}, fmt.Errorf("unsupported witness kind %q", kind)
	}
}

func ciEnvelopeProfile(kind string) bool {
	return kind == KindGitLabCI || kind == KindBuildkite
}

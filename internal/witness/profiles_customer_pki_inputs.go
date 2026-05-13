package witness

import (
	"sort"
	"strings"
)

func missingCustomerPKIInputs(opts ProfileOptions) []string {
	// Required input reporting is stable and sorted so the missing-evidence
	// surface can be compared across runs.
	missing := []string{}
	for name, value := range requiredCustomerPKIInputs(opts) {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, name)
		}
	}
	return appendMissingCustomerPKIKeyInput(missing, opts)
}

func appendMissingCustomerPKIKeyInput(missing []string, opts ProfileOptions) []string {
	// Public key and public certificate are alternatives, so the missing-input
	// message reports them as one choice.
	if strings.TrimSpace(opts.CustomerPKIPublicCert) == "" && strings.TrimSpace(opts.CustomerPKIPublicKey) == "" {
		missing = append(missing, "--customer-pki-public-cert|--customer-pki-public-key")
	}
	sort.Strings(missing)
	return missing
}

func requiredCustomerPKIInputs(opts ProfileOptions) map[string]string {
	// These three external evidence files are mandatory; the public key/cert
	// alternative is checked separately because either form can carry the key.
	return map[string]string{
		"--customer-pki-authority-policy":   opts.CustomerPKIAuthorityPolicy,
		"--customer-pki-payload-digest":     opts.CustomerPKIPayloadDigest,
		"--customer-pki-freshness-evidence": opts.CustomerPKIFreshness,
	}
}

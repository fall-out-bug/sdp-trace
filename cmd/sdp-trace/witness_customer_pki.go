package main

import (
	"sort"
	"strings"
)

func missingCustomerPKIFlags(opts *flagSet) []string {
	missing := []string{}
	// These fields establish authority, payload binding, and freshness; all are
	// required for customer-PKI witness construction.
	missing = appendMissingStringFlags(missing, opts, map[string]string{
		"customer-pki-authority-policy":   "--customer-pki-authority-policy",
		"customer-pki-payload-digest":     "--customer-pki-payload-digest",
		"customer-pki-freshness-evidence": "--customer-pki-freshness-evidence",
	})
	if missingCustomerPKIPublicCredential(opts) {
		// Either a certificate or raw public key is enough for the public
		// credential anchor.
		missing = append(missing, "--customer-pki-public-cert or --customer-pki-public-key")
	}
	// Sorted output keeps remediation deterministic.
	sort.Strings(missing)
	return missing
}

func appendMissingStringFlags(missing []string, opts *flagSet, required map[string]string) []string {
	for name, flag := range required {
		if strings.TrimSpace(opts.stringValue(name)) == "" {
			// Preserve the user-facing flag spelling in remediation output.
			missing = append(missing, flag)
		}
	}
	return missing
}

func missingCustomerPKIPublicCredential(opts *flagSet) bool {
	return strings.TrimSpace(opts.stringValue("customer-pki-public-cert")) == "" && strings.TrimSpace(opts.stringValue("customer-pki-public-key")) == ""
}

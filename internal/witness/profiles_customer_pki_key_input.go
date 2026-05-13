package witness

import (
	"crypto/ed25519"
	"os"
	"strings"
)

func privateKeyInput(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	// This best-effort preflight catches accidental private-key files before
	// parsing public-key or certificate material.
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return containsSecretLike(raw)
}

func loadCustomerPublicKey(opts ProfileOptions) (ed25519.PublicKey, error) {
	// Load the chosen public trust anchor only after input-path and private-key
	// checks have already rejected unsafe material.
	raw, err := os.ReadFile(customerPKIPublicKeyPath(opts))
	if err != nil {
		return nil, err
	}
	return parseCustomerPublicKeyPEM(raw)
}

func customerPKIPublicKeyPath(opts ProfileOptions) string {
	if opts.CustomerPKIPublicKey == "" {
		// Certificate input is the fallback trust anchor only when a direct
		// public key path was not supplied.
		return opts.CustomerPKIPublicCert
	}
	return opts.CustomerPKIPublicKey
}

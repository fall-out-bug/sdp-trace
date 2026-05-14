package main

type witnessOptions struct {
	kind                      string
	target                    string
	out                       string
	reportDir                 string
	witnessEnvelope           string
	customerPKIAuthorityPath  string
	customerPKIPublicCertPath string
	customerPKIPublicKeyPath  string
	customerPKIPayloadDigest  string
	customerPKIFreshnessPath  string
}

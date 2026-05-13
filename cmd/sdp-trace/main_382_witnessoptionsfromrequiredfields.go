package main

func witnessOptionsFromRequiredFields(fields witnessRequiredFields, opts *flagSet) witnessOptions {
	// Optional fields are copied verbatim; witness package validation decides
	// whether a specific profile can use or must reject them.
	// Required fields stay first so the resulting record always has its core
	// target, kind, and output provenance.
	return witnessOptions{
		kind:                      fields.kind,
		target:                    fields.target,
		out:                       fields.out,
		reportDir:                 opts.stringValue("report-dir"),
		witnessEnvelope:           opts.stringValue("witness-envelope"),
		customerPKIAuthorityPath:  opts.stringValue("customer-pki-authority-policy"),
		customerPKIPublicCertPath: opts.stringValue("customer-pki-public-cert"),
		customerPKIPublicKeyPath:  opts.stringValue("customer-pki-public-key"),
		customerPKIPayloadDigest:  opts.stringValue("customer-pki-payload-digest"),
		customerPKIFreshnessPath:  opts.stringValue("customer-pki-freshness-evidence"),
	}
}

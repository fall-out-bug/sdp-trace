package main

func newWitnessFlagSet() *flagSet {
	opts := &flagSet{name: "witness"}
	// Witness flags cover both generic CI witnesses and customer-PKI material;
	// semantic validation happens after parsing.
	opts.setString("kind", "")
	opts.setString("out", "")
	opts.setString("report-dir", "")
	opts.setString("witness-envelope", "")
	opts.setString("customer-pki-authority-policy", "")
	opts.setString("customer-pki-public-cert", "")
	opts.setString("customer-pki-public-key", "")
	opts.setString("customer-pki-payload-digest", "")
	opts.setString("customer-pki-freshness-evidence", "")
	return opts
}

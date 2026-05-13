package main

import (
	"fmt"
	"io"
)

func parseWitnessFlagSet(args []string, stderr io.Writer) (*flagSet, bool) {
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
	if err := opts.parse(args); err != nil {
		// Malformed witness flags stop before any CI or Customer PKI material is
		// read from disk.
		fmt.Fprintln(stderr, err)
		return nil, false
	}
	// The witness command has no positional target; target comes from flags so
	// generated records have explicit provenance fields.
	return opts, true
}

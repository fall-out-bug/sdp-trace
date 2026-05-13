package main

import (
	"strings"
)

func missingCustomerPKIPublicCredential(opts *flagSet) bool {
	return strings.TrimSpace(opts.stringValue("customer-pki-public-cert")) == "" && strings.TrimSpace(opts.stringValue("customer-pki-public-key")) == ""
}

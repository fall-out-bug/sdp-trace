package main

import (
	"fmt"
	"strings"
)

func validateWitnessKindFlags(kind string, opts *flagSet) (string, bool) {
	// Missing kind-specific material is a usage failure, not a generated
	// not_assessed witness record.
	missing := missingWitnessKindFlags(kind, opts)
	if len(missing) > 0 {
		return fmt.Sprintf("customer-pki witness requires %s", strings.Join(missing, ", ")), false
	}
	return "", true
}

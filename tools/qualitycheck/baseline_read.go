package main

import (
	"fmt"
	"go/token"
)

type schemaVersioned interface {
	schemaVersion() string
}

func readMIBaseline[T schemaVersioned](path string, wantSchema string, label string) (T, error) {
	// Read first, validate schema second: JSON shape errors and stale-ratchet
	// errors remain distinguishable in user-facing diagnostics.
	baseline, err := readJSONFile[T](path)
	if err != nil {
		return baseline, err
	}
	// Schema mismatch is a verifier failure because old ratchets may encode a
	// different subject key or MI meaning.
	if baseline.schemaVersion() != wantSchema {
		// Include the metric family label so function and file baselines fail
		// with actionable messages when both are configured.
		return baseline, fmt.Errorf("unsupported %s MI baseline schema %q", label, baseline.schemaVersion())
	}
	return baseline, nil
}

func addHalsteadToken(tok token.Token, literal string, counts *halsteadCounts) {
	// Classification is deliberately narrow: Go literals and identifiers are
	// operands, while keywords and punctuation are operators.
	if skipHalsteadToken(tok) {
		return
	}
	key := halsteadTokenKey(tok, literal)
	counts.length++
	if tok.IsLiteral() || tok == token.IDENT {
		counts.operands[key]++
		return
	}
	counts.operators[key]++
}

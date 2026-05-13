package posture

import (
	"errors"
)

var (
	errDigestMismatch  = errors.New("digest mismatch")
	errMissingRequired = errors.New("missing required input")
	errUnsafePath      = errors.New("unsafe path")
)

type digestErrorReason struct {
	err    error
	reason string
}

var digestErrorReasons = []digestErrorReason{
	{err: errDigestMismatch, reason: "untrusted_input_digest_mismatch"},
	{err: errMissingRequired, reason: "missing_required_input"},
	{err: errUnsafePath, reason: "malformed_input"},
}

// reasonForDigestErr maps digest errors to refusal reasons at the trust boundary.
// Mismatch is distinct to avoid blurring tamper evidence with cannot-verify failures.
func reasonForDigestErr(err error) string {
	// reasonForDigestErr keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	for _, item := range digestErrorReasons {
		if errors.Is(err, item.err) {
			return item.reason
		}
	}
	return "malformed_input"
}

// trustForDigestErr maps digest verification outcomes to trust states at the evidence boundary.
// Distinguishes tamper evidence (untrusted) from verification failures (cannot_verify).
func trustForDigestErr(err error) string {
	// trustForDigestErr keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if errors.Is(err, errDigestMismatch) {
		return "untrusted_input"
	}
	return "cannot_verify_input"
}

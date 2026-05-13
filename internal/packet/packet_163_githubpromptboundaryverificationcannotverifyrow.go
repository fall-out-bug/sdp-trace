package packet

import (
	"strings"
)

func githubPromptBoundaryVerificationCannotVerifyRow(required bool, classification PromptBoundaryClassification) (Row, bool) {
	// githubPromptBoundaryVerificationCannotVerifyRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if promptBoundaryBlocksVerification(required, classification) {

		return githubRow("PC-VERIFICATION", StateCannotVerify, "Verification cannot pass without clean or partially retained prompt-boundary evidence.", []string{"prompt:boundary"}, strings.Join(classification.Reasons, "; ")), true
	}
	return Row{}, false
}

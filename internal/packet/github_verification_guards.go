package packet

import "strings"

// Cannot-verify checks are separated so prompt-boundary failures are not hidden
// behind ordinary missing CI evidence.
func githubVerificationCannotVerifyRow(input GitHubPREvidenceInput, classification PromptBoundaryClassification) (Row, bool) {
	if row, ok := githubPromptBoundaryVerificationCannotVerifyRow(input.RequirePromptBoundary, classification); ok {
		return row, true
	}
	return githubCheckVerificationCannotVerifyRow(input)
}

func githubPromptBoundaryVerificationCannotVerifyRow(required bool, classification PromptBoundaryClassification) (Row, bool) {
	if promptBoundaryBlocksVerification(required, classification) {
		return githubRow("PC-VERIFICATION", StateCannotVerify, "Verification cannot pass without clean or partially retained prompt-boundary evidence.", []string{"prompt:boundary"}, strings.Join(classification.Reasons, "; ")), true
	}
	return Row{}, false
}

// CI-owned generation requires both check evidence and the workflow run id;
// otherwise a current run cannot be bound to the generated packet.
func githubCheckVerificationCannotVerifyRow(input GitHubPREvidenceInput) (Row, bool) {
	if len(input.Checks) == 0 {
		return githubRow("PC-VERIFICATION", StateCannotVerify, "No GitHub check evidence was provided.", nil, "missing GitHub check or workflow run evidence"), true
	}
	if missingRequiredWorkflowRunID(input) {
		return githubRow("PC-VERIFICATION", StateCannotVerify, "No current workflow run id was provided.", []string{"github:check"}, "missing workflow run id for CI-owned packet generation"), true
	}
	return Row{}, false
}

func missingRequiredWorkflowRunID(input GitHubPREvidenceInput) bool {
	return input.RequirePromptBoundary && strings.TrimSpace(input.WorkflowRunID) == ""
}

// Prompt-boundary fail and cannot_verify both block verification pass when a
// retained boundary is required by the generation mode.
func promptBoundaryBlocksVerification(required bool, classification PromptBoundaryClassification) bool {
	return required && (classification.RouteProofEffect == StateFail || classification.RouteProofEffect == StateCannotVerify)
}

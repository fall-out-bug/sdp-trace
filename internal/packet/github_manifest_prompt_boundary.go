package packet

import "strings"

// Prompt text takes precedence over digest metadata because retained text gives
// reviewers the strongest replay surface.
func promptBoundaryResolver(boundary PromptBoundary) string {
	if strings.TrimSpace(boundary.Text) != "" {
		return "prompt:text-retained"
	}
	if strings.TrimSpace(boundary.Digest) != "" {
		return "prompt:digest:" + boundary.Digest
	}
	return "prompt:missing"
}

// Retained-form mirrors resolver strength: text is redacted, digest-only stays
// digest_only, and absent material remains not_retained.
func promptBoundaryRetainedForm(boundary PromptBoundary) string {
	if strings.TrimSpace(boundary.Text) != "" {
		return "redacted"
	}
	if strings.TrimSpace(boundary.Digest) != "" {
		return "digest_only"
	}
	return "not_retained"
}

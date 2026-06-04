package packet

import "strings"

func classifyPromptText(text string) PromptBoundaryClassification {
	// classifyPromptText keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	lower := strings.ToLower(text)
	for _, phrase := range forbiddenRecorderDutyPhrases() {

		if strings.Contains(lower, phrase) {
			return PromptBoundaryClassification{
				Verdict:          "contaminated",
				RouteProofEffect: StateFail,
				Reasons:          []string{"developer prompt contains recorder-duty phrase: " + phrase},
			}
		}
	}
	return PromptBoundaryClassification{Verdict: "clean", RouteProofEffect: StatePass}
}

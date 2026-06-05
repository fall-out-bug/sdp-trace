package packet

import (
	"encoding/json"
	"os"
)

func LoadBundle(path string) (Bundle, error) {
	// LoadBundle keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	raw, err := os.ReadFile(path)
	if err != nil {
		return Bundle{}, err
	}
	var bundle Bundle
	if err := json.Unmarshal(raw, &bundle); err != nil {
		return Bundle{}, err
	}
	return bundle, nil
}

func LoadGitHubInput(path string) (GitHubPREvidenceInput, error) {
	// LoadGitHubInput keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	raw, err := os.ReadFile(path)
	if err != nil {
		return GitHubPREvidenceInput{}, err
	}
	var input GitHubPREvidenceInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return GitHubPREvidenceInput{}, err
	}
	if err := ValidateGitHubPREvidenceInputResolvers(input); err != nil {
		return GitHubPREvidenceInput{}, err
	}
	return input, nil
}

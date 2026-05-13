package packet

func checksSucceeded(checks []GitHubCheck) bool {
	// checksSucceeded keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, check := range checks {
		if check.Conclusion != "success" {

			return false
		}
	}
	return true
}

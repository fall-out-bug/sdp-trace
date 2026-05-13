package repoobserver

import "strings"

const (
	// Managed markers bound the only .gitignore region install mode may edit.
	gitignoreBeginMarker = "# sdp-trace begin"
	gitignoreEndMarker   = "# sdp-trace end"
	gitignoreBlock       = "# sdp-trace begin\n.sdp-trace/hooks/\n.sdp-trace/ci/\n.sdp-trace/install-diff.txt\n# sdp-trace end\n"
)

func locateGitignoreBlock(text string) (int, int) {
	// Both markers must be ordered correctly before the block is considered
	// manageable.
	start := strings.Index(text, gitignoreBeginMarker)
	if start < 0 {
		return -1, -1
	}
	end := strings.Index(text, gitignoreEndMarker)
	if end < start {
		return -1, -1
	}
	// Include the end marker so replacement preserves surrounding user content.
	return start, end + len(gitignoreEndMarker)
}

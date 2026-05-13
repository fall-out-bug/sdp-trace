package repoobserver

import (
	"fmt"
	"os"
	"strings"
)

func updateSdpTraceGitignoreBlock(opts Options, path, text string, data []byte, start, end int) ([]DiffSummary, error) {
	// Replacing an existing managed block requires --force and a backup because
	// the user may have edited it.
	// Exact generated block matches are idempotent and produce no summary.
	// Replacement preserves all text before and after the managed marker range.
	// The backup stores the full previous .gitignore, not just the managed block.
	// Diff summaries again use digest/size metadata rather than raw ignore
	// content.
	// Force mode is the only path that can replace a divergent managed block.
	current := text[start:end]
	if current == strings.TrimSuffix(gitignoreBlock, "\n") {
		return nil, nil
	}
	if !opts.Force {
		return nil, fmt.Errorf("%s: .gitignore sdp-trace block differs; use --force after reviewing safe diff", ReasonManualStepRequired)
	}
	if err := os.WriteFile(path+".bak", data, 0o644); err != nil {
		return nil, fmt.Errorf("%s: backup failed for .gitignore", ReasonUnsafeOutputRefused)
	}
	next := text[:start] + strings.TrimSuffix(gitignoreBlock, "\n") + text[end:]
	return replacedGitignoreBlockSummary(data, next), os.WriteFile(path, []byte(next), 0o644)
}

func replacedGitignoreBlockSummary(before []byte, next string) []DiffSummary {
	// The summary is digest and count metadata only; the managed ignore content
	// is not copied into force-mode review output.
	return []DiffSummary{{
		Path:    ".gitignore",
		Action:  "replace_sdp_trace_block",
		Before:  contentSummary(before),
		After:   contentSummary([]byte(next)),
		Summary: "replace marked sdp-trace gitignore block using safe byte and line counts",
		Backup:  ".gitignore.bak",
	}}
}

func appendGitignoreMarker(path, text string) ([]DiffSummary, error) {
	// Preserve existing ignore content and append a clean newline boundary before
	// the managed block.
	if strings.TrimSpace(text) != "" && !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	text += gitignoreBlock
	return nil, os.WriteFile(path, []byte(text), 0o644)
}

package repoobserver

import "strings"

func safeRef(ref string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	// Normalize path separators before redacting unsafe config values in
	// human-facing output.
	ref = strings.ReplaceAll(ref, "\\", "/")
	if strings.HasPrefix(ref, "/") || strings.Contains(ref, ":/") {
		return "unsafe_absolute_path_redacted"
	}
	// Relative refs stay actionable in reports while absolute paths are
	// redacted above.
	return ref
}

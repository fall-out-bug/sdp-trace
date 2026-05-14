package main

// homePathAllowlist records historical specs that contain absolute local paths.
// Each entry is grandfathered because the file is historical evidence that has
// already been reviewed; new .md files must not add absolute paths.
var homePathAllowlist = map[string]bool{
	"specs/004-mvp-readiness-hardening/completion-audit.md":                          true,
	"specs/004-mvp-readiness-hardening/implementation-ledger.md":                     true,
	"specs/004-mvp-readiness-hardening/pr-review/ec8db52/packet/inputs/context-2.md": true,
	"specs/004-mvp-readiness-hardening/pr-review/ec8db52/packet/inputs/context-3.md": true,
	"specs/007-github-oss-demo-packet/reviews/2026-05-10-socratic-review.md":         true,
	"specs/008-invisible-flight-recorder/reviews/raw/post-code-deepseek.md":          true,
	"specs/008-invisible-flight-recorder/reviews/raw/post-dx-deepseek.md":            true,
	"specs/008-invisible-flight-recorder/reviews/raw/post-evidence-deepseek.md":      true,
}

package repoobserver

import (
	"fmt"
	"strings"
)

func writeHumanTableHeader(b *strings.Builder, status Status) {
	// The repository root is represented abstractly; absolute checkout paths are
	// not part of the portable status table.
	fmt.Fprintf(b, "Profile: %s\n", status.Profile)
	fmt.Fprintf(b, "Repository: %s\n", status.RepositoryID)
	fmt.Fprintf(b, "Install state: %s\n", status.InstallState)
	fmt.Fprintf(b, "Proof state: %s\n\n", status.ProofState)
	b.WriteString("Surface | Install state | Proof state | Trust scope | Evidence source | Next action\n")
	b.WriteString("--- | --- | --- | --- | --- | ---\n")
}

func writeHumanTableSurfaces(b *strings.Builder, surfaces []Surface) {
	for _, surface := range surfaces {
		// Render every surface independently so install and proof gaps remain
		// inspectable instead of being collapsed into a health score.
		writeHumanTableSurface(b, surface)
	}
}

func writeHumanTableSurface(b *strings.Builder, surface Surface) {
	// Empty remediation renders as "-" so absence of an action is explicit.
	// The table prints install and proof states side by side to avoid implying
	// that installed means verified.
	action := surface.NextAction
	if action == "" {
		action = "-"
	}
	fmt.Fprintf(b, "%s | %s | %s | %s | %s | %s\n",
		surface.SurfaceID,
		surface.InstallState,
		surface.ProofState,
		surface.TrustScope,
		surface.EvidenceSource,
		action,
	)
}

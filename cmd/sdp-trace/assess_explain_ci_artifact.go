package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
)

func explainCIArtifactObservation(result ciartifact.ObservationResult, stdout io.Writer) int {
	// CI artifact explanations keep producer, access, binding, and output
	// safety states separate so one gap cannot mask another.
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "CI artifact observation: %s\n", result.ArtifactObservationState)
	// Top-level scopes summarize the overall observation; family rows below
	// retain their own finer-grained state.
	fmt.Fprintf(stdout, "Authority scope: %s\n", result.AuthorityScope)
	fmt.Fprintf(stdout, "Producer scope: %s\n", result.ProducerScope)
	fmt.Fprintf(stdout, "Artifact access state: %s\n", result.ArtifactAccessState)
	explainCIArtifactFamilies(result.ArtifactFamilies, stdout)
	// Index and output-safety summaries are printed after family rows so they
	// cannot be confused with per-family evidence.
	fmt.Fprintf(stdout, "Artifact index: %s (%s)\n", result.ArtifactIndex.Result, result.ArtifactIndex.ReasonCode)
	fmt.Fprintf(stdout, "Output safety: %s (%s)\n", result.OutputSafety.State, result.OutputSafety.ReasonCode)
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
	return 0
}

func explainCIArtifactFamilies(families []ciartifact.FamilyObservation, stdout io.Writer) {
	for _, family := range families {
		// Each family line preserves producer, access, and binding state as
		// separate evidence dimensions.
		fmt.Fprintf(stdout, "Artifact family %s: %s (%s)\n", family.Family, family.FamilyState, family.ReasonCode)
		fmt.Fprintf(stdout, "  Producer scope: %s\n", family.ProducerScope)
		fmt.Fprintf(stdout, "  Artifact access: %s\n", family.ArtifactAccessState)
		fmt.Fprintf(stdout, "  Binding: %s\n", family.BindingState)
	}
}

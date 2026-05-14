package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
)

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

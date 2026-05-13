package verifier

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func readExplainManifest(runDir string) (trace.RunManifest, error) {
	// The manifest supplies stable display context only; it is not treated as a
	// verdict source here.
	manifestPath := filepath.Join(runDir, "run.json")
	var manifest trace.RunManifest
	if err := trace.ReadJSON(context.Background(), manifestPath, &manifest); err != nil {
		return trace.RunManifest{}, fmt.Errorf("run directory missing: %w", err)
	}
	return manifest, nil
}

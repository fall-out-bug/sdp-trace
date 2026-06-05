package prreview

import (
	"os"
	"path/filepath"
	"time"
)

func RunReview(packet Packet, profile ReviewProfile, opts RunOptions) (RunSet, *RunPreview, error) {
	// RunReview validates or projects review data; it does not create external proof.
	// Profile validation intentionally happens before option normalization or
	// filesystem preparation so malformed review policy cannot create artifacts.
	if err := validateProfile(profile); err != nil {
		return RunSet{}, nil, err
	}
	opts = normalizeRunOptions(opts)
	if opts.Preview {
		return RunSet{}, preview(packet, profile), nil
	}
	return runReview(packet, profile, opts)
}

// normalizeRunOptions fills replay metadata defaults while preserving explicit
// caller-provided values for deterministic harness runs.
func normalizeRunOptions(opts RunOptions) RunOptions {
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}
	if opts.WorkDir == "" {
		opts.WorkDir = "."
	}
	return opts
}

// prepareRunDirectories owns the creation boundary for run artifacts: the run
// directory must be new, while the raw subdirectory is materialized inside it.
func prepareRunDirectories(outDir string) (string, error) {
	if err := ensureNewDir(outDir); err != nil {
		return "", err
	}
	rawDir := filepath.Join(outDir, "raw")
	if err := os.MkdirAll(rawDir, 0o755); err != nil {
		return "", err
	}
	return rawDir, nil
}

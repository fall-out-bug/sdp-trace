package harnessobs

import "fmt"

func resolveObservePaths(opts ObserveOptions) (string, string, string, error) {
	profilePath, err := safeExistingFile(opts.ProfilePath)
	if err != nil {
		return "", "", "", fmt.Errorf("unsafe profile path: %w", err)
	}
	sourcePath, err := safeExistingFile(opts.SourcePath)
	if err != nil {
		return "", "", "", fmt.Errorf("unsafe source path: %w", err)
	}
	outDir, err := safeOutDir(opts.OutDir)
	if err != nil {
		return "", "", "", err
	}
	return profilePath, sourcePath, outDir, nil
}

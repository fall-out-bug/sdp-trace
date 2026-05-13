package demo

func discoverRunDirsUnder(root string) ([]string, error) {
	// discoverRunDirsUnder keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if err := ensureRunRootDir(root); err != nil {
		return nil, err
	}
	if hasRunManifest(root) {
		return []string{root}, nil
	}
	return collectRunDirs(root)
}

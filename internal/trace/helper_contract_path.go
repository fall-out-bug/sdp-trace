package trace

import "path/filepath"

// ResolveContractPath converts empty or relative paths to absolute file system paths.
func ResolveContractPath(baseDir, contractPath string) string {
	// Empty paths stay empty because the default contract may be selected without
	// a filesystem artifact.
	if contractPath == "" {
		return ""
	}
	if filepath.IsAbs(contractPath) {
		return contractPath
	}
	if baseDir == "" {
		return contractPath
	}
	return filepath.Join(baseDir, contractPath)
}

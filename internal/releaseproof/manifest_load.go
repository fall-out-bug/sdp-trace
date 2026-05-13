package releaseproof

import (
	"crypto/sha256"
	"encoding/json"
	"os"
)

func loadManifest(repoRoot, manifestPath string) (manifestData, error) {
	// Resolve the manifest through the repository boundary before hashing so
	// the verification record names the same relative file that was read.
	manifestRel, manifestBytes, err := loadManifestBytes(repoRoot, manifestPath)
	if err != nil {
		return manifestData{}, err
	}
	var manifest Manifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return manifestData{}, err
	}
	// The digest is computed over the exact bytes accepted for parsing.
	return manifestData{
		manifest: manifest,
		ref:      manifestRel,
		digest:   sha256.Sum256(manifestBytes),
	}, nil
}

func loadManifestBytes(repoRoot, manifestPath string) (string, []byte, error) {
	// Normalize and resolve before reading so ManifestRef and ManifestDigest
	// describe the same repository-contained file.
	manifestRel, err := cleanRepoRelativePath(manifestPath)
	if err != nil {
		return "", nil, err
	}
	manifestAbs, err := resolveRepoFile(repoRoot, manifestRel)
	if err != nil {
		return "", nil, err
	}
	manifestBytes, err := os.ReadFile(manifestAbs)
	return manifestRel, manifestBytes, err
}

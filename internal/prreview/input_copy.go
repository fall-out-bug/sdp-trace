package prreview

import (
	"os"
	"path/filepath"
)

// Input copying materializes external evidence into the packet directory.
//
// The copy path preserves source bytes, produces stable digest refs, and makes
// file permissions explicit after write so local umask cannot weaken the packet
// contract observed by downstream verifiers.
func copyInputs(inputDir, prefix string, paths []string) ([]SafeRef, error) {
	refs := make([]SafeRef, 0, len(paths))
	for i, path := range paths {
		ref, err := copiedInputRef(inputDir, prefix, i, path)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

func copyInput(inputDir, name, source, kind, contentType string) (SafeRef, error) {
	data, err := os.ReadFile(source)
	if err != nil {
		return SafeRef{}, err
	}
	if err := writeCopiedInput(inputDir, name, data); err != nil {
		return SafeRef{}, err
	}
	return copiedInputSafeRef(name, kind, contentType, data), nil
}

func writeCopiedInput(inputDir, name string, data []byte) error {
	dest := filepath.Join(inputDir, name)
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return err
	}
	return os.Chmod(dest, 0o644)
}

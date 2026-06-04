package prreview

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
)

// Copied input refs define the public packet shape.
//
// Names are position-based, extensions are normalized through the package
// content-type rules, and refs are always packet-relative paths under inputs/.
// Callers may specialize Kind after copying for context or verification inputs.
func copiedInputRef(inputDir, prefix string, index int, path string) (SafeRef, error) {
	position := index + 1
	name := fmt.Sprintf("%s-%d%s", prefix, position, normalizedExt(path))
	ref, err := copyInput(inputDir, name, path, RefKindDoc, contentType(path))
	if err != nil {
		return SafeRef{}, err
	}
	ref.ID = fmt.Sprintf("%s-%d", prefix, position)
	return ref, nil
}

func copiedInputSafeRef(name, kind, contentType string, data []byte) SafeRef {
	digest := sha256.Sum256(data)
	return SafeRef{
		ID:             strings.TrimSuffix(name, filepath.Ext(name)),
		Kind:           kind,
		Ref:            filepath.ToSlash(filepath.Join("inputs", name)),
		DigestSHA256:   hex.EncodeToString(digest[:]),
		ContentType:    contentType,
		RedactionState: RedactionNone,
	}
}

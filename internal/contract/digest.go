package contract

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

// Digest computes a deterministic digest for fixture signing and lock checks.
func (c ExpectedEvidenceContract) Digest() (string, error) {
	// Canonicalization sorts unordered contract slices before hashing so the
	// digest describes contract meaning rather than authoring order.
	canonical, err := canonicalizeContract(c)
	if err != nil {
		return "", err
	}

	return canonicalDigest(canonical), nil
}

func canonicalizeContract(contract ExpectedEvidenceContract) ([]byte, error) {
	// Work on a value copy so digest generation never mutates caller-owned
	// slices that may be reused by tests or later validation.
	stable := contract

	sort.Strings(stable.RequiredObservers)
	sort.Strings(stable.OptionalObservers)
	sort.Strings(stable.RequiredEvents)
	sort.Strings(stable.GateEvents)
	canonical, err := trace.CanonicalJSON(stable)
	if err != nil {
		return nil, err
	}
	return canonical, nil
}

func canonicalDigest(data []byte) string {
	// Lowercase hex is the digest representation used across trace fixtures.
	sum := sha256.Sum256(data)
	return strings.ToLower(hex.EncodeToString(sum[:]))
}

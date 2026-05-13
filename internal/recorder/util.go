package recorder

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

// Utility helpers are kept small and deterministic because they feed artifact
// naming, manifest digests, and run identifiers used by verifier replay.
//
// Helpers in this file intentionally avoid recorder policy decisions. They only
// adapt standard library behavior into the narrow shapes needed by manifests,
// events, and artifact files.
//
// That separation keeps trust rules in the recorder workflow and leaves these
// routines reusable by tests without hidden side effects.
//
// When a helper must be lossy, the surrounding comment names the fallback so it
// is not mistaken for proof-grade evidence.
//
// All persistent writes still return errors to their callers.
//
// No helper here emits a gate verdict.

func commandName(command []string) string {
	// Empty command slices are validated by the CLI layer; this helper remains
	// defensive for tests and direct package callers.
	if len(command) == 0 {
		return ""
	}

	return filepath.Base(command[0])
}

func currentWorkingDir() string {
	// The recorder treats cwd lookup failures as an empty label because source
	// snapshotting already carries the stronger provenance signal.
	cwd, _ := os.Getwd()
	return cwd
}

func mustMarshalJSON(value any) []byte {
	// This helper is used only for digest material where returning a stable empty
	// object is preferable to panicking inside recorder setup.
	data, err := json.Marshal(value)
	if err != nil {
		return []byte("{}")
	}
	return data
}

func writeIndentedJSON(path string, value any) error {
	// Human-readable JSON keeps run artifacts useful during audits while still
	// leaving event hashes to canonical event computation.
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func writeText(path string, value string) error {
	// Text artifacts are intentionally plain files so shell and verifier tooling
	// can compare digests without a JSON parser.
	return os.WriteFile(path, []byte(value), 0o644)
}

func ifNilCopy(values []string) []string {
	// Nil env means "no inherited environment"; non-nil env is copied so command
	// preparation does not alias caller-owned slices.
	if values == nil {
		return []string{}
	}
	return append([]string{}, values...)
}

func randomHex(length int) string {
	// Recorder identifiers are opaque labels, not security tokens, but using
	// crypto randomness avoids predictable collisions in local evidence trees.
	const alphabet = "0123456789abcdef"
	out := make([]byte, length)
	for i := range out {
		v, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return fmt.Sprintf("fallback-%x", time.Now().UnixNano())
		}
		out[i] = alphabet[v.Int64()]
	}
	return string(out)
}

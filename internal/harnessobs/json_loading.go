package harnessobs

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
)

// JSON loading helpers share the same safe existing-file boundary while
// preserving the caller's choice between permissive and strict decoding.
//
// Permissive loading is used for already-versioned artifacts where unknown
// fields are tolerated by the caller. Strict loading is used for author-owned
// profile inputs where unknown fields and trailing JSON must be rejected before
// later validation reasons are trusted.

// readExistingJSON loads a safe local file and leaves unknown-field handling to
// the target type and standard JSON unmarshaler.
func readExistingJSON(path string, target any) error {
	safePath, err := safeExistingFile(path)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(safePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// readExistingJSONStrict loads a safe local file and rejects payload shape that
// would otherwise be ignored by ordinary unmarshaling.
func readExistingJSONStrict(path string, target any) error {
	safePath, err := safeExistingFile(path)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(safePath)
	if err != nil {
		return err
	}
	return decodeStrictJSON(data, target)
}

// decodeStrictJSON rejects unknown fields and additional JSON values after the
// first object so profile inputs have a single authoritative payload.
func decodeStrictJSON(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}

	// A second decode distinguishes trailing JSON from ordinary EOF.
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("strict JSON input contains trailing data")
	}
	return nil
}

package prreview

import (
	"os"
	"path/filepath"
)

// ReadPacket accepts either a packet JSON path or a packet directory containing
// the canonical packet.json artifact.
func ReadPacket(path string) (Packet, error) {
	var packet Packet
	if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {
		path = filepath.Join(path, "packet.json")
	}
	return packet, readJSON(path, &packet)
}

// ReadProfile decodes a review profile and immediately validates its schema
// contract so invalid profiles cannot enter review orchestration.
func ReadProfile(path string) (ReviewProfile, error) {
	var profile ReviewProfile
	if err := readJSON(path, &profile); err != nil {
		return profile, err
	}
	return profile, validateProfile(profile)
}

// ReadRunSet accepts either results.json or a run output directory, then
// validates run identity and per-result contracts after decode.
func ReadRunSet(path string) (RunSet, error) {
	var runs RunSet
	path = runSetPath(path)
	if err := readJSON(path, &runs); err != nil {
		return runs, err
	}
	if err := validateRunSet(runs); err != nil {
		return runs, err
	}
	return runs, nil
}

// runSetPath resolves the canonical results.json path when callers pass a run
// output directory instead of the file itself.
func runSetPath(path string) string {
	if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {
		return filepath.Join(path, "results.json")
	}
	return path
}

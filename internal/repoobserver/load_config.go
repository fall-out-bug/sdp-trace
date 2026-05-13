package repoobserver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadConfig(repoRoot string) (Config, error) {
	// Config is local structural input; malformed config is unsafe because it
	// could mislabel repository identity or profile.
	data, err := os.ReadFile(filepath.Join(repoRoot, ".sdp-trace", "config.json"))
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("%s: .sdp-trace/config.json is malformed", ReasonUnsafeOutputRefused)
	}
	return validateConfig(config)
}

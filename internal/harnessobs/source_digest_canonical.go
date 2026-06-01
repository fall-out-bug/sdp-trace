package harnessobs

import "encoding/json"

func canonicalSourceDigestLine(line []byte) ([]byte, bool) {
	var raw map[string]any
	if err := json.Unmarshal(line, &raw); err != nil {
		return nil, false
	}
	raw["source_digest"] = ""

	canonical, err := json.Marshal(raw)
	if err != nil {
		return nil, false
	}
	return canonical, true
}

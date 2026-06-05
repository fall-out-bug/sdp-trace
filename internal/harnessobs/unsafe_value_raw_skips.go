package harnessobs

import "strings"

// skippableRawEventField allows known raw body/tool prompt positions that are
// consumed during normalization but must not become retained evidence fields.
func skippableRawEventField(path, key string, value any, rawEvent bool) bool {
	return rawEvent &&
		(unretainedRawToolInputField(path, key, value) ||
			(unretainedRawBodyField(key) && !structuredRawBody(value)))
}

// unretainedRawToolInputField recognizes the OpenCode tool input prompt path
// that may appear in raw input but is excluded from retained event evidence.
func unretainedRawToolInputField(path, key string, value any) bool {
	if key != "prompt" {
		return false
	}
	if _, ok := value.(string); !ok {
		return false
	}

	segments := strings.Split(path, ".")
	if len(segments) < 3 {
		return false
	}
	return path == "part.state.input.prompt"
}

// unretainedRawBodyField lists free-text raw body fields that may be skipped
// only when the value is not structured retained evidence.
func unretainedRawBodyField(key string) bool {
	switch key {
	case "text", "content", "input", "output", "stdout", "stderr":
		return true
	default:
		return false
	}
}

// structuredRawBody prevents map-shaped raw bodies from being skipped as opaque
// text; nested values still need recursive unsafe checks.
func structuredRawBody(value any) bool {
	switch value.(type) {
	case map[string]any:
		return true
	default:
		return false
	}
}

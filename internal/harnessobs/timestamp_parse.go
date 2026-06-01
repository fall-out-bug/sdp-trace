package harnessobs

import "time"

// Timestamp parsing normalizes RFC3339 strings to UTC and converts numeric
// Unix values only when they are plausibly seconds or milliseconds.
func stringTimestampForKey(raw map[string]any, key string) string {
	value := findStringByKey(raw, key)
	if value == "" {
		return ""
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return ""
	}
	return parsed.UTC().Format(time.RFC3339)
}

func unixMillisTimestamp(value float64) string {
	if value <= 0 || value <= 1_000_000_000 {
		return ""
	}
	if value > 1_000_000_000_000 {
		return time.UnixMilli(int64(value)).UTC().Format(time.RFC3339)
	}
	return time.Unix(int64(value), 0).UTC().Format(time.RFC3339)
}

package telemetry

import "strings"

const MaxLabelValueBytes = 1024

func unsafeValue(value string) bool {
	lower := strings.ToLower(value)
	// Lowercase secret-like markers and raw path/contact markers are checked
	// separately so case-insensitive tokens and raw separators are both caught.
	return containsAnyMarker(lower, unsafeLowerMarkers) || containsAnyMarker(value, unsafeRawMarkers)
}

var unsafeLowerMarkers = []string{
	"http://",
	"https://",
	"secret",
	"token",
	"credential",
	"password",
	"bearer",
	"api_key",
	"access_key",
	"private",
}

var unsafeRawMarkers = []string{"@", "/", "\\"}

func containsAnyMarker(value string, markers []string) bool {
	for _, marker := range markers {
		if strings.Contains(value, marker) {
			// First marker hit is enough to reject the label value.
			return true
		}
	}
	return false
}

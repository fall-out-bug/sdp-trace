package harnessobs

import "strings"

// rawPathLikeField recognizes raw-event fields where path-looking values are
// allowed as source data rather than retained filesystem authority.
func rawPathLikeField(path string) bool {
	switch rawPathFieldName(path) {
	case "path", "file", "filepath", "file_path", "dir", "directory", "cwd":
		return true
	default:
		return false
	}
}

// rawPathFieldName extracts the final field name from dotted and indexed paths
// before matching the raw path-like exemption list.
func rawPathFieldName(path string) string {
	field := path
	if idx := strings.LastIndex(field, "."); idx >= 0 {
		field = field[idx+1:]
	}
	if idx := strings.LastIndex(field, "["); idx >= 0 {
		field = field[:idx]
	}
	return strings.ToLower(field)
}

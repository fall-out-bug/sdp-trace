package harnessobs

func unsupportedRawEventFormat(format string) bool {
	return format != "" && format != OpenCodeJSONLRawFormat
}

package query

func sensitiveClasses() []string {
	classes := append([]string(nil), verifiedAbsentSensitiveClasses...)
	return append(classes, verifiedAbsentProviderClasses...)
}

package adaptercapture

func validRetentionMode(mode string) bool {
	return validRetentionModes[mode]
}

var validRetentionModes = map[string]bool{
	RetentionDigestOnly:          true,
	RetentionSanitizedExcerpt:    true,
	RetentionEncryptedRawRef:     true,
	RetentionExternalArtifactRef: true,
	RetentionNotAssessed:         true,
}

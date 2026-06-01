package harnessobs

import "strings"

// Safe token rendering keeps generated identifiers bounded and path-safe while
// preserving the historical default token for empty values.
// Each accepted rune is one byte from `safeTokenRunes`, so truncating after
// rendering preserves the old 128-byte cap without splitting UTF-8.
func safeToken(value string) string {
	token := strings.Trim(safeTokenBody(value), "-_.:")
	if token == "" {
		return "opencode"
	}
	return token
}

func safeTokenBody(value string) string {
	var b strings.Builder
	for _, r := range value {
		writeSafeTokenRune(&b, r)
	}
	return truncateSafeToken(b.String())
}

func truncateSafeToken(token string) string {
	if len(token) > 128 {
		return token[:128]
	}
	return token
}

func writeSafeTokenRune(b *strings.Builder, r rune) {
	if safeTokenRune(r) {
		b.WriteRune(r)
		return
	}
	b.WriteByte('-')
}

func safeTokenRune(r rune) bool {
	return strings.ContainsRune(safeTokenRunes, r)
}

package witness

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

func parseOIDCClaims(token string) (OIDCClaims, error) {
	// Claim parsing intentionally reads only the JWT payload; signature
	// validation is outside this portable witness helper.
	payload, err := decodeJWTPayload(token)
	if err != nil {
		return OIDCClaims{}, err
	}
	return oidcClaimsFromPayload(payload)
}

func decodeJWTPayload(token string) ([]byte, error) {
	// Only the payload segment is decoded; the raw JWT is not retained in any
	// output artifact.
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, errors.New("invalid jwt shape")
	}
	return base64.RawURLEncoding.DecodeString(parts[1])
}

func oidcClaimsFromPayload(payload []byte) (OIDCClaims, error) {
	var raw rawOIDCClaims
	// Decode into the raw JWT schema first so audience normalization stays
	// isolated from JSON parsing.
	if err := json.Unmarshal(payload, &raw); err != nil {
		return OIDCClaims{}, err
	}
	return normalizedOIDCClaims(raw), nil
}

func normalizedOIDCClaims(raw rawOIDCClaims) OIDCClaims {
	// Only claims needed for environment binding are retained; unused JWT
	// material is discarded at the trust boundary.
	return OIDCClaims{
		Issuer:     raw.Issuer,
		Subject:    raw.Subject,
		Audience:   audienceString(raw.Audience),
		Repository: raw.Repository,
		Ref:        raw.Ref,
		SHA:        raw.SHA,
	}
}

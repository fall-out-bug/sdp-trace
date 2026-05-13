package witness

// OIDCClaims keeps only the GitHub token claims required to compare the live
// identity token with the environment snapshot.
type OIDCClaims struct {
	Issuer     string `json:"issuer"`
	Subject    string `json:"subject"`
	Audience   string `json:"audience"`
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	SHA        string `json:"sha"`
}

// rawOIDCClaims mirrors the JWT payload shape before audience normalization.
type rawOIDCClaims struct {
	Issuer     string `json:"iss"`
	Subject    string `json:"sub"`
	Audience   any    `json:"aud"`
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	SHA        string `json:"sha"`
}

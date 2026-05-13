package witness

import (
	"net/http"
)

const githubOIDCIssuer = "https://token.actions.githubusercontent.com"

type TokenFetcher func(env map[string]string) (string, error)

func FetchGitHubOIDCToken(env map[string]string) (string, error) {
	// The live fetcher keeps request construction centralized so tests can
	// inject a deterministic TokenFetcher without touching network behavior.
	token, err := fetchGitHubOIDCToken(
		http.DefaultClient,
		env["ACTIONS_ID_TOKEN_REQUEST_URL"],
		env["ACTIONS_ID_TOKEN_REQUEST_TOKEN"],
	)
	if err != nil {
		return "", err
	}
	return token, nil
}

func fetchGitHubOIDCToken(httpClient *http.Client, requestURL, requestToken string) (string, error) {
	// Request creation, transport, and response parsing are split so each trust
	// boundary can fail without exposing the bearer token.
	req, err := buildOIDCTokenRequest(requestURL, requestToken)
	if err != nil {
		return "", err
	}
	body, err := executeOIDCTokenRequest(httpClient, req)
	if err != nil {
		return "", err
	}
	return parseOIDCTokenResponse(body)
}

package witness

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func executeOIDCTokenRequest(httpClient *http.Client, req *http.Request) ([]byte, error) {
	// Response bodies are returned only for successful status codes so provider
	// error payloads are not copied into witness diagnostics.
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !successHTTPStatus(resp.StatusCode) {
		return nil, fmt.Errorf("oidc token request returned %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func successHTTPStatus(statusCode int) bool {
	// Treat only HTTP 2xx as token response success.
	return statusCode >= 200 && statusCode < 300
}

func parseOIDCTokenResponse(body []byte) (string, error) {
	// The GitHub OIDC endpoint returns the token in `value`; no other response
	// field is retained by witness generation.
	var payload struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}
	if payload.Value == "" {
		// Empty token responses cannot establish identity and are not replaced
		// with environment-only evidence.
		return "", errors.New("oidc token response missing value")
	}
	return payload.Value, nil
}

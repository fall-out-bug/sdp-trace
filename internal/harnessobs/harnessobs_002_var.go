package harnessobs

import (
	"errors"

	"regexp"
)

var (
	errSessionSourceUnavailable = errors.New("session source unavailable")
	safeIDPattern               = regexp.MustCompile(`^[A-Za-z0-9_.:-]{1,128}$`)
	safeFileIDPattern           = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)
	sha256Pattern               = regexp.MustCompile(`^[a-f0-9]{64}$`)
	base64TokenPattern          = regexp.MustCompile(`(?i)^[A-Za-z0-9+/_-]{43,}={0,2}$`)
	providerTokenPrefix         = regexp.MustCompile(`(?i)(sk-[A-Za-z0-9_-]{16,}|xox[baprs]-[A-Za-z0-9-]{16,}|gh[pousr]_[A-Za-z0-9_]{20,})`)
	privatePathPattern          = regexp.MustCompile(`(^|[\s"'])/(Users|home|private|var|tmp)/[^\s"']+`)
	rawFieldNames               = map[string]bool{
		"raw_prompt":         true,
		"prompt":             true,
		"raw_model_response": true,
		"model_response":     true,
		"raw_command":        true,
		"command_body":       true,
	}
	sensitiveFieldNames = map[string]bool{
		"access_token":  true,
		"api_key":       true,
		"apikey":        true,
		"authorization": true,
		"auth":          true,
		"token":         true,
	}
	authQueryKeys = map[string]bool{
		"token": true, "access_token": true, "api_key": true, "apikey": true,
		"key": true, "signature": true, "sig": true, "auth": true,
	}
)

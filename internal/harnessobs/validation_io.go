package harnessobs

import (
	"encoding/json"
	"io"
)

func nonAuthority() string {
	return "harness observation is evidence only; no harness compliance, feature delivery, PR approval, merge approval, release readiness, or production trust is claimed"
}

func DecodeValidation(r io.Reader) (Validation, error) {
	var validation Validation
	if err := json.NewDecoder(r).Decode(&validation); err != nil {
		return Validation{}, err
	}
	return validation, nil
}

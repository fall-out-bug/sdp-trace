package harnessobs

import "encoding/json"

func validationDigest(validation Validation) string {
	copy := validation
	copy.ValidationDigest = ""

	data, _ := json.Marshal(copy)
	return sha256Hex(data)
}

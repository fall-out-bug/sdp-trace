package harnessobs

import "encoding/json"

func digestCommand(command []string) string {
	data, _ := json.Marshal(command)
	return sha256Hex(data)
}

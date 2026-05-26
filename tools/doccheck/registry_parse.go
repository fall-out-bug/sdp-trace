package main

import (
	"encoding/json"
	"fmt"
)

func parseCommandSurface(output []byte) (commandSurface, error) {
	var surface commandSurface
	if err := json.Unmarshal(output, &surface); err != nil {
		return commandSurface{}, fmt.Errorf("parse command surface: %w", err)
	}
	return surface, nil
}

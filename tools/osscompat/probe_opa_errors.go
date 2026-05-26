package main

import "fmt"

func opaJSONError(err error) error {
	return fmt.Errorf("opa eval output is not valid JSON: %w", err)
}

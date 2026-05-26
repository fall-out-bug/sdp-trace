package main

import "fmt"

func opaExpressionBool(value interface{}) (bool, error) {
	v, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("opa eval result is not a boolean")
	}
	return v, nil
}

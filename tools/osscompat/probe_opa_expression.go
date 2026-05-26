package main

import "fmt"

func checkOPAExpressions(result []opaExpressionSet) error {
	if len(result) == 0 || len(result[0].Expressions) == 0 {
		return fmt.Errorf("opa eval returned no expressions")
	}
	return nil
}

package main

import "fmt"

func assertOPAResult(pass bool, expectPass bool, label string) (verifierState, string) {
	if pass == expectPass {
		if expectPass {
			// Positive fixtures must evaluate to true.
			return statePass, fmt.Sprintf("adapter.rego evaluates %s as expected", label)
		}
		// Negative fixtures pass the probe only when policy rejects them.
		return statePass, fmt.Sprintf("adapter.rego correctly rejects %s", label)
	}
	if expectPass {
		return stateFail, "opa eval did not return boolean true for expected pass fixture"
	}
	return stateFail, fmt.Sprintf("opa eval did not return boolean false for %s negative fixture", label)
}

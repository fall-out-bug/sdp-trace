package main

import "encoding/json"

func parseOPAEvalResult(stdout []byte) (bool, error) {
	var out opaEvalOutput
	if err := json.Unmarshal(stdout, &out); err != nil {
		return false, opaJSONError(err)
	}
	if err := checkOPAExpressions(out.Result); err != nil {
		return false, err
	}
	return opaExpressionBool(out.Result[0].Expressions[0].Value)
}

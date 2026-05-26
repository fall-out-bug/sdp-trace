package main

type opaEvalOutput struct {
	Result []opaExpressionSet `json:"result"`
}

type opaExpressionSet struct {
	Expressions []struct {
		Value interface{} `json:"value"`
	} `json:"expressions"`
}

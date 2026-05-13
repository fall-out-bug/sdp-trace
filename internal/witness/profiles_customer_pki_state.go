package witness

var customerPKIStateSetters = map[string]func(*ProfileStates){
	"artifact":  func(states *ProfileStates) { states.ArtifactBindingState = stateFail },
	"freshness": func(states *ProfileStates) { states.FreshnessState = stateFail },
	"policy":    func(states *ProfileStates) { states.PolicyBindingState = stateFail },
	"run":       func(states *ProfileStates) { states.RunBindingState = stateFail },
	"signer":    func(states *ProfileStates) { states.SignerAuthorityState = stateFail },
}

type profileDecision struct {
	status string
	scope  string
	reason string
}

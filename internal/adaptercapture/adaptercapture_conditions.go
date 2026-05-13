package adaptercapture

func pass(id, code, reason string) Condition {
	return Condition{ID: id, State: StatePass, ReasonCode: code, Reason: reason}
}

func fail(id, code, reason, action string) Condition {
	return Condition{ID: id, State: StateFail, ReasonCode: code, Reason: reason, NextAction: action}
}

func cannotVerify(id, code, reason, action string) Condition {
	return Condition{ID: id, State: StateCannotVerify, ReasonCode: code, Reason: reason, NextAction: action}
}

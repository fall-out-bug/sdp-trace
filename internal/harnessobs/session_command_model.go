package harnessobs

import "strings"

func hasSessionCommandModel(session SessionRun) bool {
	return session.CommandModelState == StatePass && strings.TrimSpace(session.CommandModel) != ""
}

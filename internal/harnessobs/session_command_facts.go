package harnessobs

import "encoding/json"

func sessionCommandFacts(session SessionRun) []Event {
	if !hasSessionCommandModel(session) {
		return nil
	}
	event := sessionCommandModelEvent(session)
	data, err := json.Marshal(event)
	if err != nil {
		return nil
	}
	event.SourceDigest = digestLine(data)
	return []Event{event}
}

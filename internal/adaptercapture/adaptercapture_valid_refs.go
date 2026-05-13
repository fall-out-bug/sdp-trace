package adaptercapture

func validProviderRefs(source string) []ProviderRef {
	return []ProviderRef{{SourceRef: "repo:generic/source", SourceCommit: source, ChangeRef: "change:42", ReviewRef: "review:7", Producer: "generic_git_host", ObservedAt: "2026-05-07T10:00:00Z"}}
}

func validEventFamilySummaries() []EventFamilyState {
	return []EventFamilyState{{EventFamily: "tool_call", State: StatePass, RetentionMode: RetentionSanitizedExcerpt, Reconstructable: true}}
}

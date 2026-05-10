package query

import "testing"

func TestMapSourceState(t *testing.T) {
	t.Parallel()

	cases := []struct {
		state string
		want  string
	}{
		{state: "pass", want: RowStatePresent},
		{state: "", want: RowStateCannotVerify},
		{state: "fail", want: RowStateIssueObserved},
		{state: RowStateCannotVerify, want: RowStateCannotVerify},
		{state: RowStateNotAssessed, want: RowStateNotAssessed},
		{state: RowStateMissingTelemetry, want: RowStateMissingTelemetry},
		{state: RowStateUnsupported, want: RowStateUnsupported},
		{state: RowStateNotIntegrated, want: RowStateNotIntegrated},
		{state: RowStateRetentionLimited, want: RowStateRetentionLimited},
		{state: "unknown", want: RowStateCannotVerify},
	}

	for _, tc := range cases {
		if got := mapSourceState(tc.state); got != tc.want {
			t.Fatalf("mapSourceState(%q) = %q, want %q", tc.state, got, tc.want)
		}
	}
}

package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

type assessCLICase[T any] struct {
	name      string
	mutate    func(*T)
	wantExit  int
	wantState string
	wantCode  string
}

func runAssessCLICases[T any](t *testing.T, cases []assessCLICase[T], valid func() T, inputName, stateField string, args func(root, inputPath string) []string, assertNoLeak func(*testing.T, string)) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			input := valid()
			tc.mutate(&input)
			inputPath := filepath.Join(root, inputName)
			writeTestJSON(t, inputPath, input)
			var out bytes.Buffer
			var errOut bytes.Buffer
			exit := run(args(root, inputPath), &out, &errOut)
			if exit != tc.wantExit {
				t.Fatalf("exit = %d want %d err=%s out=%s", exit, tc.wantExit, errOut.String(), out.String())
			}
			if !strings.Contains(out.String(), `"`+stateField+`": "`+tc.wantState+`"`) ||
				!strings.Contains(out.String(), `"reason_code": "`+tc.wantCode+`"`) {
				t.Fatalf("output missing state/code: %s", out.String())
			}
			assertNoLeak(t, out.String())
		})
	}
}

package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestFlagSetParsesQuotedStringValue(t *testing.T) {
	flags := &flagSet{name: "wrap"}
	flags.setString("name", "")

	if err := flags.parse([]string{"--name", "agent run payload", "--", "echo", "with space"}); err != nil {
		t.Fatalf("parse returned error: %v", err)
	}
	if got, want := flags.stringValue("name"), "agent run payload"; got != want {
		t.Fatalf("name = %q, want %q", got, want)
	}
	if got, want := flags.rest(), []string{"echo", "with space"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("rest = %#v, want %#v", got, want)
	}
}

func TestFlagSetRejectsMissingStringValueAtEnd(t *testing.T) {
	flags := &flagSet{name: "wrap"}
	flags.setString("name", "")

	err := flags.parse([]string{"--name"})
	if err == nil {
		t.Fatalf("expected parse error")
	}
	if !strings.Contains(err.Error(), "requires value") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFlagSetRejectsUnknownFlagWithEquals(t *testing.T) {
	flags := &flagSet{name: "wrap"}
	flags.setString("name", "")

	err := flags.parse([]string{"--outpt-dir=path"})
	if err == nil {
		t.Fatalf("expected parse error")
	}
	if !strings.Contains(err.Error(), "unknown flag --outpt-dir") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFlagSetParsesBooleanValues(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		want      bool
		wantErr   bool
		errSubstr string
	}{
		{name: "implicit true", args: []string{"--enabled"}, want: true},
		{name: "explicit true", args: []string{"--enabled=true"}, want: true},
		{name: "explicit one", args: []string{"--enabled=1"}, want: true},
		{name: "explicit false", args: []string{"--enabled=false"}, want: false},
		{name: "explicit zero", args: []string{"--enabled=0"}, want: false},
		{name: "invalid literal", args: []string{"--enabled=maybe"}, wantErr: true, errSubstr: "invalid boolean value"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			flags := &flagSet{name: "wrap"}
			flags.setBool("enabled", false)

			err := flags.parse(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected parse error")
				}
				if !strings.Contains(err.Error(), tt.errSubstr) {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("parse returned error: %v", err)
			}
			if got := flags.boolValue("enabled"); got != tt.want {
				t.Fatalf("enabled = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlagSetRepeatsFlagsOverwriteAndKeepOrder(t *testing.T) {
	flags := &flagSet{name: "wrap"}
	flags.setString("name", "")
	flags.setBool("enabled", true)

	err := flags.parse([]string{
		"--name", "first",
		"--name", "second",
		"--enabled=false",
		"--enabled",
		"echo", "done",
	})
	if err != nil {
		t.Fatalf("parse returned error: %v", err)
	}
	if got, want := flags.stringValue("name"), "second"; got != want {
		t.Fatalf("name = %q, want %q", got, want)
	}
	if got, want := flags.boolValue("enabled"), true; got != want {
		t.Fatalf("enabled = %v, want %v", got, want)
	}
	if got, want := flags.rest(), []string{"echo", "done"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("rest = %#v, want %#v", got, want)
	}
}

func TestFlagSetCapturesEverythingAfterSeparator(t *testing.T) {
	flags := &flagSet{name: "wrap"}
	flags.setString("name", "")

	err := flags.parse([]string{"--name", "demo", "--", "--name", "ignored", "--bad", "x"})
	if err != nil {
		t.Fatalf("parse returned error: %v", err)
	}
	if got, want := flags.rest(), []string{"--name", "ignored", "--bad", "x"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("rest = %#v, want %#v", got, want)
	}
}

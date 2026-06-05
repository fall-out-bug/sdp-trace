package harnessobs

import (
	"encoding/json"
	"testing"
)

func TestLoadProfileValidatesProfile(t *testing.T) {
	dir := t.TempDir()
	path := "profile.json"
	oldwd := chdir(t, dir)
	defer oldwd()

	writeJSONFixture(t, path, validProfileFixture())
	if _, err := LoadProfile(path); err != nil {
		t.Fatalf("LoadProfile() error = %v", err)
	}
	assertProfileRuleJSONShape(t)

	profile := validProfileFixture()
	profile.ProfileID = "../bad"
	writeJSONFixture(t, path, profile)
	if _, err := LoadProfile(path); err == nil || err.Error() != "unsafe profile_id" {
		t.Fatalf("LoadProfile() validation error = %v, want unsafe profile_id", err)
	}
}

func assertProfileRuleJSONShape(t *testing.T) {
	t.Helper()
	profile := Profile{DegradationRules: map[string]Rule{
		"empty_rule":     {},
		"populated_rule": {State: StateNotAssessed, ReasonCode: "required_event_family_absent"},
	}}
	data, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("marshal profile rule shape: %v", err)
	}
	var rawProfile map[string]any
	if err := json.Unmarshal(data, &rawProfile); err != nil {
		t.Fatalf("parse profile rule shape: %v", err)
	}
	rawRules, ok := rawProfile["degradation_rules"].(map[string]any)
	if !ok {
		t.Fatalf("degradation_rules = %#v, want object", rawProfile["degradation_rules"])
	}
	assertRawRule(t, rawRules, "empty_rule", "", "")
	assertRawRule(t, rawRules, "populated_rule", StateNotAssessed, "required_event_family_absent")
}

func assertRawRule(t *testing.T, rawRules map[string]any, key, wantState, wantReasonCode string) {
	t.Helper()
	rawRule, ok := rawRules[key].(map[string]any)
	if !ok {
		t.Fatalf("degradation_rules[%q] = %#v, want object", key, rawRules[key])
	}
	wantRule := map[string]any{
		"state":       wantState,
		"reason_code": wantReasonCode,
	}
	for field, want := range wantRule {
		if got := rawRule[field]; got != want {
			t.Fatalf("degradation_rules[%q][%q] = %#v, want %#v in %#v", key, field, got, want, rawRule)
		}
	}
}

func TestValidateProfile(t *testing.T) {
	tests := []struct {
		name    string
		profile Profile
		wantErr string
	}{
		{
			name:    "passes with valid profile",
			profile: validProfileFixture(),
			wantErr: "",
		},
		{
			name: "unsupported schema version",
			profile: func() Profile {
				p := validProfileFixture()
				p.SchemaVersion = "bad"
				return p
			}(),
			wantErr: "unsupported harness profile schema_version: bad",
		},
		{
			name: "unsafe profile id",
			profile: func() Profile {
				p := validProfileFixture()
				p.ProfileID = "../bad"
				return p
			}(),
			wantErr: "unsafe profile_id",
		},
		{
			name: "unsafe harness family",
			profile: func() Profile {
				p := validProfileFixture()
				p.HarnessFamily = "bad family"
				return p
			}(),
			wantErr: "unsafe harness_family",
		},
		{
			name: "unsupported event schema version",
			profile: func() Profile {
				p := validProfileFixture()
				p.EventSchemaVersion = "bad"
				return p
			}(),
			wantErr: "unsupported event_schema_version",
		},
		{
			name: "missing required family",
			profile: func() Profile {
				p := validProfileFixture()
				p.RequiredEventFamilies = nil
				return p
			}(),
			wantErr: "profile requires at least one required_event_family",
		},
		{
			name: "unsupported required family",
			profile: func() Profile {
				p := validProfileFixture()
				p.RequiredEventFamilies = []string{"harness", "bad-family"}
				return p
			}(),
			wantErr: "unsupported event family: bad-family",
		},
		{
			name: "unsupported optional family",
			profile: func() Profile {
				p := validProfileFixture()
				p.OptionalEventFamilies = []string{"bad-family"}
				return p
			}(),
			wantErr: "unsupported event family: bad-family",
		},
		{
			name: "unsupported degradation rule",
			profile: func() Profile {
				p := validProfileFixture()
				p.DegradationRules = map[string]Rule{
					"bad-key": {State: StatePass, ReasonCode: "ok"},
				}
				return p
			}(),
			wantErr: "unsupported degradation rule: bad-key",
		},
		{
			name: "invalid degradation rule state",
			profile: func() Profile {
				p := validProfileFixture()
				p.DegradationRules = map[string]Rule{
					"missing_required_family": {State: "bad", ReasonCode: "bad"},
				}
				return p
			}(),
			wantErr: "invalid degradation rule missing_required_family",
		},
		{
			name: "invalid degradation rule reason code",
			profile: func() Profile {
				p := validProfileFixture()
				p.DegradationRules = map[string]Rule{
					"missing_required_family": {State: StateNotAssessed, ReasonCode: "bad reason"},
				}
				return p
			}(),
			wantErr: "invalid degradation rule missing_required_family",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProfile(tt.profile)
			if err == nil {
				if tt.wantErr != "" {
					t.Fatalf("validateProfile() = nil, want %q", tt.wantErr)
				}
				return
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("validateProfile() = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func validProfileFixture() Profile {
	return Profile{
		SchemaVersion:         ProfileSchemaVersion,
		ProfileID:             "generic-harness-v1",
		HarnessFamily:         "generic-harness",
		EventSchemaVersion:    EventSchemaVersion,
		RequiredEventFamilies: []string{"harness"},
		OptionalEventFamilies: []string{"model"},
		RawRetentionPolicy:    "digest_only",
		DegradationRules: map[string]Rule{
			"missing_required_family": {State: StateNotAssessed, ReasonCode: "required_event_family_absent"},
			"missing_optional_family": {State: StateNotAssessed, ReasonCode: "optional_event_family_absent"},
			"source_unavailable":      {State: StateCannotVerify, ReasonCode: "source_unavailable"},
			"unsafe_input":            {State: StateFail, ReasonCode: "unsafe_input"},
			"digest_mismatch":         {State: StateCannotVerify, ReasonCode: "source_digest_mismatch"},
			"schema_version_mismatch": {State: StateCannotVerify, ReasonCode: "schema_version_mismatch"},
			"cross_link_conflict":     {State: StateCannotVerify, ReasonCode: "adapter_harness_state_conflict"},
		},
	}
}

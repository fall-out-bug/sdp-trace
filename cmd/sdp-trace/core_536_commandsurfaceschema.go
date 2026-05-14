package main

// commandSurfaceSchema defines the machine-readable command surface.
// This is the source of truth for help/docs drift checks and agent discovery.
type commandSurfaceSchema struct {
	SchemaVersion string              `json:"schema_version"`
	Commands      []commandSurfaceCmd `json:"commands"`
	Profiles      []profileMeta       `json:"profiles,omitempty"`
	WitnessKinds  []string            `json:"witness_kinds,omitempty"`
	States        []stateMeta         `json:"states,omitempty"`
}

// commandSurfaceCmd describes one top-level command family.
// State may be "complete", "partial", or "not_assessed".
type commandSurfaceCmd struct {
	Name          string           `json:"name"`
	Usage         string           `json:"usage"`
	Variations    []string         `json:"variations,omitempty"`
	Description   string           `json:"description,omitempty"`
	Subcommands   []string         `json:"subcommands,omitempty"`
	RequiredFlags []flagMeta       `json:"required_flags,omitempty"`
	OptionalFlags []flagMeta       `json:"optional_flags,omitempty"`
	RepeatedFlags []flagMeta       `json:"repeated_flags,omitempty"`
	Positional    []positionalMeta `json:"positional,omitempty"`
	RestBehavior  string           `json:"rest_behavior,omitempty"`
	OutputPaths   []outputPathMeta `json:"output_paths,omitempty"`
	TrustNote     string           `json:"trust_note,omitempty"`
	State         string           `json:"state"`
}

// flagMeta describes a CLI flag.
type flagMeta struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Description  string `json:"description,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
}

// positionalMeta describes a positional argument.
type positionalMeta struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
}

// outputPathMeta describes where output is written.
type outputPathMeta struct {
	Flag        string `json:"flag,omitempty"`
	Positional  int    `json:"positional,omitempty"`
	Description string `json:"description,omitempty"`
}

// profileMeta describes an assessment profile.
type profileMeta struct {
	ID          string `json:"id"`
	Description string `json:"description,omitempty"`
	Command     string `json:"command,omitempty"`
}

// stateMeta describes a result or trust state value.
type stateMeta struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

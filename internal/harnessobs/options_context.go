package harnessobs

import "time"

type ObserveOptions struct {
	ProfilePath string
	SourcePath  string
	OutDir      string
	Now         time.Time
}

type SessionSetupOptions struct {
	ProfilePath string
	OutDir      string
	Command     string
	Now         time.Time
}

type SessionCollectOptions struct {
	ProfilePath string
	RunDir      string
	Now         time.Time
}

type SessionOptions struct {
	ProfilePath string
	OutDir      string
	Command     []string
	Now         time.Time
}

type ValidateOptions struct {
	ProfilePath string
	RunDir      string
	OutPath     string
}

type observationContext struct {
	outDir       string
	sourcePath   string
	sourceDigest string
	now          time.Time
	profile      Profile
	events       []Event
}

type sessionCollectionContext struct {
	profilePath        string
	runDir             string
	now                time.Time
	profile            SessionProfile
	session            SessionRun
	harnessProfile     Profile
	harnessProfilePath string
}

type observedCommandResult struct {
	waitErr error
	end     time.Time
}

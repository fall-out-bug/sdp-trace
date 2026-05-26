package main

import "runtime"

// envInfo captures the benchmark environment.
type envInfo struct {
	Platform  string `json:"platform"`
	GoOS      string `json:"goos"`
	GoArch    string `json:"goarch"`
	GoVersion string `json:"go_version"`
}

func getEnv() envInfo {
	return envInfo{
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
		GoOS:      runtime.GOOS,
		GoArch:    runtime.GOARCH,
		GoVersion: runtime.Version(),
	}
}

package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/authority"
	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
)

var ciArtifactExitCodes = map[string]int{
	ciartifact.StatePass: 0,
	ciartifact.StateFail: 1,
}

var authorityExitCodes = map[string]int{
	authority.StateWithinAuthority:  0,
	authority.StateOutsideAuthority: 1,
	authority.StateNotAssessed:      exitCannotVerify,
}

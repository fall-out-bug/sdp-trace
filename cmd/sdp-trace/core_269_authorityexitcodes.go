package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/authority"
)

var authorityExitCodes = map[string]int{
	authority.StateWithinAuthority:  0,
	authority.StateOutsideAuthority: 1,
	authority.StateNotAssessed:      exitCannotVerify,
}

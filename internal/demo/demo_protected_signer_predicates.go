package demo

import "github.com/fall_out_bug/sdp-trace/internal/checkpoint"

func protectedSignerPass(state, trustScope string) bool {
	// CI-signed checkpoints are the only passing protected signer path.
	return state == GatePass && trustScope == checkpoint.TrustScopeCISigned
}

func protectedSignerLocalOnly(state, trustScope string) bool {
	// Local signatures remain useful evidence but cannot satisfy protected mode.
	return state == GatePass && trustScope == checkpoint.TrustScopeLocalSigned
}

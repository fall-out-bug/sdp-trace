package repoobserver

func repositoryIdentitySurface(opts Options) Surface {
	if opts.RepositoryID != "" {
		// Caller-supplied identity binds the observation to an intended source,
		// but it is still local structural evidence rather than external proof.
		return surface(SurfaceRepositoryIdentity, StatePass, StateNotAssessed, ScopeLocalStructural, "caller_supplied_repository_id", ReasonManualStepRequired, "", "")
	}
	// A sanitized origin hash avoids leaking remotes while keeping the identity
	// surface inspectable for follow-up binding work.
	return surface(SurfaceRepositoryIdentity, StatePass, StateNotAssessed, ScopeLocalStructural, "sanitized_origin_hash", ReasonManualStepRequired, "", "")
}

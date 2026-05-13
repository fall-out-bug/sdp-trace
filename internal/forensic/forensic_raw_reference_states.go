package forensic

func rawReferenceAccessUnverifiable(ref *RawReference) bool {
	// rawReferenceAccessUnverifiable keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	switch ref.AccessState {
	case "", AccessStateNotAssessed, AccessStateUnavailable, AccessStateRevoked:

		return true
	default:
		return false
	}
}

func encryptedKeyCustodyUnverifiable(ref *RawReference) bool {
	// encryptedKeyCustodyUnverifiable keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	if ref.ReferenceType != RetentionModeEncryptedRawRef {

		return false
	}
	switch ref.KeyCustodyState {
	case "", KeyCustodyUnknown, KeyCustodyNotAssessed, KeyCustodyCompromised, KeyCustodyDestroyed:

		return true
	default:
		return false
	}
}

func retentionLifecycleUnverifiable(state string) bool {
	// retentionLifecycleUnverifiable keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	switch state {
	case "", RetentionLifecycleNotAssessed, RetentionLifecycleExpired, RetentionLifecycleRevoked, RetentionLifecycleDeleted:

		return true
	default:
		return false
	}
}

package adaptercapture

func adapterEventOverclaims(event AdapterEvent) bool {
	return event.ReconstructableClaimed && adapterEventInsufficient(event) && event.CapAnnotation == ""
}

func adapterEventInsufficient(event AdapterEvent) bool {
	return event.CaptureState != "captured" || insufficientRetentionModes[event.RetentionMode]
}

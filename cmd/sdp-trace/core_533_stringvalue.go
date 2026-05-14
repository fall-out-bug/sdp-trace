package main

func (f *flagSet) stringValue(key string) string {
	if f.data == nil {
		// Unregistered string maps read as absent flags, matching parse defaults.
		return ""
	}
	return f.data[key]
}

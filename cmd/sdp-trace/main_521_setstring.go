package main

func (f *flagSet) setString(key, defaultValue string) {
	if f.data == nil {
		// Lazily allocate so tiny commands only register the flags they own.
		f.data = map[string]string{}
	}
	f.data[key] = defaultValue
}

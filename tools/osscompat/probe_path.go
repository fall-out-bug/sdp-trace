package main

import "os"

func requireExistingPath(path string) error {
	_, err := os.Stat(path)
	return err
}

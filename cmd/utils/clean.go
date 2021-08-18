package utils

import "os"

func CleanUp(path string) error {
	if CheckIfPathExists(path) {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}
	return nil
}

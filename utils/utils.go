package utils

import (
	"os"
)

func MakeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModeDir|0755)
	}
	return nil
}

func CheckIfPathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CleanUp(path string) error {
	if CheckIfPathExists(path) {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}
	return nil
}

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

func CheckIfInArray(arr []string, val string) bool {
	for _, i := range arr {
		if i == val {
			return true
		}
	}
	return false
}

func CheckIfMultipleInArray(arr []string, val string) bool {
	counter := 0
	for _, i := range arr {
		if i == val {
			counter++
		}
		if counter > 1 {
			return true
		}
	}
	return false
}

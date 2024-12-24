package utils

import (
	"path/filepath"
	"runtime"
)

func RemoveEmptyStringInSlice(s []string) []string {
	var newS []string
	for _, a := range s {
		if a != "" {
			newS = append(newS, a)
		}
	}
	return newS
}

func GetCurrentFileName() string {
	_, filename, _, ok := runtime.Caller(2)
	if !ok {
		panic("No caller information")
	}
	return filename
}

func GetCurrentFolderName() string {
	return filepath.Dir(GetCurrentFileName())
}

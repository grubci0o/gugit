package util

import (
	"path/filepath"
	"strings"
)

func Ignore(path string) bool {
	path = filepath.ToSlash(path)

	if strings.HasPrefix(filepath.Base(path), ".") {
		return true
	}

	ignorePaths := []string{
		".ugit",
		".git",
		"vendor",
		"gugit.exe",
		"main.exe",
	}

	parts := strings.Split(path, "/")
	for _, part := range parts {
		for _, ignore := range ignorePaths {
			if strings.EqualFold(part, ignore) {
				return true
			}
		}
	}

	return false
}

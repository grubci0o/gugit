package util

import (
	"path/filepath"
	"strings"
)

func Ignore(path string) bool {
	// Normalize path to forward slashes
	path = filepath.ToSlash(path)

	// Check if path starts with . or contains /.
	if strings.HasPrefix(filepath.Base(path), ".") {
		return true
	}

	// Additional specific ignores
	ignorePaths := []string{
		".ugit",
		".git",
		"vendor",
		"gugit.exe",
		"main.exe",
	}

	// Check if path or any parent directory should be ignored
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

package cmd

import (
	"gugit/internal/memory"
	"gugit/internal/util"
	"log"
)

func createTag(name, oid string) {
	updateRef("refs/tags/"+name, RefValue{
		symbolic: false,
		value:    oid,
	}, true)
}

func Tag(name, oid string) error {
	if oid == "" {
		_, val, err := getRef("HEAD", true)
		if err != nil {
			return &util.GuGitError{
				Type:    util.ErrRef,
				Message: "Failed to get HEAD reference",
				Err:     err,
			}
		}
		createTag(name, val.value)
	} else {
		createTag(name, oid)
	}
	return nil
}

// ResolveName resolves various types of references to their commit OIDs.
// Supported formats:
// - @ or HEAD: current HEAD
// - refs/heads/<branch>: branch reference
// - refs/tags/<tag>: tag reference
// - <sha1>: direct commit OID
// - <branch>: short branch name
// - <tag>: short tag name
func ResolveName(name string) string {
	// Handle special cases
	switch name {
	case "@", "HEAD":
		_, rVal, err := getRef("HEAD", true)
		if err != nil || rVal.value == "" {
			log.Println("HEAD not found or empty")
			return ""
		}
		return rVal.value
	}

	// Try to resolve in order of precedence
	lookupPaths := []string{
		name,                 // Direct reference
		"refs/" + name,       // Full reference
		"refs/tags/" + name,  // Tag reference
		"refs/heads/" + name, // Branch reference
	}

	for _, path := range lookupPaths {
		_, ref, err := getRef(path, false)
		if err == nil && ref.value != "" {
			// Dereference the ref if it exists
			_, resolvedRef, err := getRef(path, true)
			if err == nil {
				return resolvedRef.value
			}
		}
	}

	// Check if it's a valid OID
	if memory.ValidOID(name) {
		return name
	}

	return ""
}

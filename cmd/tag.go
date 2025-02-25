package cmd

import (
	"gugit/internal/memory"
	"log"
)

func createTag(name, oid string) {
	updateRef("refs/tags/"+name, RefValue{
		symbolic: false,
		value:    oid,
	}, true)
}

func Tag(name, oid string) {
	if oid == "" {
		_, val := getRef("HEAD", true)
		createTag(name, val.value)
	} else {
		createTag(name, oid)
	}
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
		exists, rVal := getRef("HEAD", true)
		if exists == "false" || rVal.value == "" {
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
		exists, ref := getRef(path, false)
		if exists != "false" && ref.value != "" {
			// Dereference the ref if it exists
			exists, resolvedRef := getRef(path, true)
			if exists != "false" {
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

package cmd

import (
	"fmt"
	"log"
	"strings"
)

func Checkout(name string) {
	// First resolve the name to get the commit OID
	oid := ResolveName(name)
	if oid == "" {
		log.Fatalf("Could not resolve reference: %s", name)
	}

	// Get the commit and read its tree
	com := getCommit(oid)
	ReadTree(com.Tree)

	// Update HEAD differently based on whether we're checking out a branch or commit
	if isBranch(name) {
		// When checking out a branch, point HEAD to the branch ref
		fmt.Printf("Switched to branch '%s'\n", strings.TrimPrefix(name, "refs/heads/"))
		updateRef("HEAD", RefValue{
			symbolic: true,
			value:    "refs/heads/" + name,
		}, false)
	} else {
		// When checking out a commit, point HEAD directly to the commit
		fmt.Printf("HEAD is now at %s\n", oid[:8])
		updateRef("HEAD", RefValue{
			symbolic: false,
			value:    oid,
		}, false)
	}
}

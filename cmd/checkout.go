package cmd

import (
	"fmt"
	"log"
	"strings"
)

func Checkout(name string) {
	oid := ResolveName(name)
	if oid == "" {
		log.Fatalf("Could not resolve reference: %s", name)
	}

	com := getCommit(oid)
	ReadTree(com.Tree)

	if isBranch(name) {
		fmt.Printf("Switched to branch '%s'\n", strings.TrimPrefix(name, "refs/heads/"))
		updateRef("HEAD", RefValue{
			symbolic: true,
			value:    "refs/heads/" + name,
		}, false)
	} else {
		fmt.Printf("HEAD is now at %s\n", oid[:8])
		updateRef("HEAD", RefValue{
			symbolic: false,
			value:    oid,
		}, false)
	}
}
